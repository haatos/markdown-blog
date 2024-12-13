package handler

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/haatos/markdown-blog/internal"
	"github.com/haatos/markdown-blog/internal/data"
	"github.com/haatos/markdown-blog/internal/model"
	"github.com/haatos/markdown-blog/internal/templates"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

var googleOAuthConfig *oauth2.Config
var githubOAuthConfig *oauth2.Config

func newGoogleOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  internal.Settings.BaseURL() + "/auth/oauth/google/callback",
		ClientID:     internal.Settings.GoogleClientID,
		ClientSecret: internal.Settings.GoogleClientSecret,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

func newGithubOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  internal.Settings.BaseURL() + "/auth/oauth/github/callback",
		ClientID:     githubOAuthConfig.ClientID,
		ClientSecret: githubOAuthConfig.ClientSecret,
		Scopes:       []string{"read:user", "user:email"},
		Endpoint:     github.Endpoint,
	}
}

const (
	oauthGoogleUserInfoURL   = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="
	oauthGithubUserURL       = "https://api.github.com/user"
	oauthGithubUserEmailsURL = "https://api.github.com/user/emails"
)

func (h *Handler) GetOAuthFlow(c echo.Context) error {
	provider := c.Param("provider")
	switch provider {
	case "google":
		if googleOAuthConfig == nil {
			googleOAuthConfig = newGoogleOAuthConfig()
		}
		oauthState := h.generateOAuthStateCookie(c)
		url := googleOAuthConfig.AuthCodeURL(oauthState)
		return c.Redirect(http.StatusTemporaryRedirect, url)
	case "github":
		if githubOAuthConfig == nil {
			githubOAuthConfig = newGithubOAuthConfig()
		}
		oauthState := h.generateOAuthStateCookie(c)
		url := githubOAuthConfig.AuthCodeURL(oauthState)
		return c.Redirect(http.StatusTemporaryRedirect, url)
	default:
		return c.Render(http.StatusUnprocessableEntity, "toast-error", templates.NewError("Invalid OAuth provider"))
	}
}

func (h *Handler) GetOAuthCallback(c echo.Context) error {
	provider := c.Param("provider")
	var user model.User
	var redirectPath string
	var err error
	switch provider {
	case "google":
		user, redirectPath, err = h.handleGoogleCallback(c)
	case "github":
		user, redirectPath, err = h.handleGithubCallback(c)
	}
	if err != nil {
		slog.Error("err handling google oauth callback", "err", err)
		return c.Redirect(http.StatusSeeOther, "/login")
	}

	if err := data.WithTx(h.rwdb, func(tx *sql.Tx) error {
		if err := data.ReadUserByEmail(c.Request().Context(), tx, &user); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				if errCreate := data.CreateUser(c.Request().Context(), tx, &user); errCreate != nil {
					return err
				}
			}
		}
		return nil
	}); err != nil {
		c.Logger().Error(err)
		return h.renderStatusPage(c, http.StatusInternalServerError)
	}

	sess := model.Session{
		ID:      generateRandomSessionID(),
		UserID:  user.ID,
		Expires: time.Now().UTC().Add(internal.Settings.SessionExpires),
	}
	if err := data.CreateSession(c.Request().Context(), h.rwdb, &sess); err != nil {
		c.Logger().Error(err)
		return h.renderStatusPage(c, http.StatusInternalServerError)
	}

	value := map[string]any{}
	value["session_id"] = sess.ID
	if err := setSessionCookie(c, value); err != nil {
		c.Logger().Error(err)
		return h.renderStatusPage(c, http.StatusInternalServerError)
	}

	return c.Redirect(http.StatusSeeOther, redirectPath)
}

func (h *Handler) generateOAuthStateCookie(c echo.Context) string {
	var expiration = time.Now().Add(1 * time.Hour)
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{
		Name:    internal.OAuthCookie,
		Value:   state,
		Expires: expiration,
	}
	c.SetCookie(&cookie)

	return state
}

func (h *Handler) handleGoogleCallback(c echo.Context) (model.User, string, error) {
	u := model.User{}
	path := "/"
	oauthState, _ := c.Cookie(internal.OAuthCookie)
	defer setCookie(c, internal.OAuthCookie, "", time.Now().UTC())

	if c.FormValue("state") != oauthState.Value {
		c.Logger().Error("invalid oauth google state")
		return u, path, errors.New("invalid OAuth Google state")
	}

	user, err := h.getGoogleUserData(c.FormValue("code"))
	if err != nil {
		return u, path, fmt.Errorf("err getting user data from google: %+v", err)
	}

	u.FirstName = user.FirstName
	u.LastName = user.LastName
	u.Email = user.Email
	u.AvatarURL = user.AvatarURL
	u.RoleID = internal.Member

	return u, path, nil
}

func (h *Handler) getGoogleUserData(code string) (model.GoogleUser, error) {
	gu := model.GoogleUser{}
	token, err := googleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		return gu, err
	}
	response, err := http.Get(oauthGoogleUserInfoURL + token.AccessToken)
	if err != nil {
		return gu, err
	}
	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&gu)
	return gu, err
}

func (h *Handler) handleGithubCallback(c echo.Context) (model.User, string, error) {
	u := model.User{}
	path := "/"
	var err error
	oauthState, _ := c.Cookie(internal.OAuthCookie)
	defer setCookie(c, internal.OAuthCookie, "", time.Now().UTC())

	if c.FormValue("state") != oauthState.Value {
		c.Logger().Error("invalid oauth github state")
		return u, path, errors.New("invalid OAuth Github state")
	}

	code := c.QueryParam("code")
	token, err := githubOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		return u, path, err
	}
	gu, err := getGithubUser(token.AccessToken)
	if err != nil {
		return u, path, err
	}

	u.Email = gu.Email
	u.FirstName, u.LastName = gu.GetFirstNameLastName()
	u.AvatarURL = gu.AvatarURL
	u.RoleID = internal.Member

	return u, path, nil
}

func getGithubUser(accessToken string) (model.GithubUser, error) {
	gu, err := getGithubUserData(accessToken)
	if err != nil {
		return gu, err
	}

	if gu.Email == "" {
		// if github user's email is not set public, we must retrieve
		// the email from another endpoint
		githubEmail, err := getUserEmailFromGithub(accessToken)
		if err != nil {
			return gu, err
		}
		gu.Email = githubEmail
	}

	return gu, nil
}

func getGithubUserData(accessToken string) (model.GithubUser, error) {
	gu := model.GithubUser{}

	req, err := http.NewRequest(http.MethodGet, oauthGithubUserURL, nil)
	if err != nil {
		return gu, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("token %s", accessToken))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return gu, err
	}

	if err := json.NewDecoder(res.Body).Decode(&gu); err != nil {
		return gu, err
	}
	return gu, nil
}

func getUserEmailFromGithub(accessToken string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, oauthGithubUserEmailsURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", fmt.Sprintf("token %s", accessToken))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	responseEmails := []struct {
		// contains other fields as well, but these
		// are the only one's we're interested in
		Email   string `json:"email"`
		Primary bool   `json:"primary"`
	}{}

	if err := json.NewDecoder(res.Body).Decode(&responseEmails); err != nil {
		return "", err
	}

	for _, re := range responseEmails {
		if re.Primary {
			return re.Email, nil
		}
	}

	return "", errors.New("no primary email found")
}
