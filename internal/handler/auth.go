package handler

import (
	"net/http"
	"time"

	"github.com/haatos/markdown-blog/internal"
	"github.com/haatos/markdown-blog/internal/templates"
	"github.com/labstack/echo/v4"
)

type LogInForm struct {
	GoogleOAuthHref string
	GithubOAuthHref string
	Email           string
	EmailError      string
}

type SignUpForm struct {
	GoogleOAuthHref string
	GithubOAuthHref string
	FirstName       string
	FirstNameError  string
	LastName        string
	LastNameError   string
	Email           string
	EmailError      string
}

type LogInPage struct {
	templates.Page
	Form LogInForm
}

func (h *Handler) GetLoginPage(c echo.Context) error {
	page := LogInPage{
		Page: getDefaultPage(c),
		Form: LogInForm{
			GoogleOAuthHref: "/auth/oauth/google",
			GithubOAuthHref: "/auth/oauth/github",
		},
	}

	tn := "log-in"
	if isHXRequest(c) {
		tn += "-main"
	}
	return c.Render(http.StatusOK, tn, page)
}

type SignUpPage struct {
	templates.Page
	Form SignUpForm
}

func (h *Handler) GetSignUpPage(c echo.Context) error {
	page := SignUpPage{
		Page: getDefaultPage(c),
		Form: SignUpForm{
			GoogleOAuthHref: "/auth/oauth/google",
			GithubOAuthHref: "/auth/oauth/github",
		},
	}

	tn := "sign-up"
	if isHXRequest(c) {
		tn += "-main"
	}
	return c.Render(http.StatusOK, tn, page)
}

func (h *Handler) GetLogOut(c echo.Context) error {
	setCookie(c, internal.SessionCookie, "", time.Now().UTC())
	return c.Redirect(http.StatusSeeOther, "/")
}
