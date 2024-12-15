package handler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/haatos/markdown-blog/internal"
	"github.com/haatos/markdown-blog/internal/data"
	"github.com/haatos/markdown-blog/internal/model"
	"github.com/labstack/echo/v4"
)

func (h *Handler) SessionMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var user *model.User

		// retrieve user from cookie & database if path is not a static path or the favicon
		r := c.Request()
		if !strings.HasPrefix(r.URL.Path, "/static") && !strings.HasPrefix(r.URL.Path, "/images") && !strings.HasPrefix(r.URL.Path, "/favicon.ico") {
			var err error
			user, err = getUserFromSessionCookie(r, h.rdb)
			if err != nil {
				c.Logger().Error("err getting user from session cookie: ", err)
				setCookie(c, internal.SessionCookie, "", time.Now().UTC())
				user = nil
			}
		}

		c.Set("user", user)
		return next(c)
	}
}

func getUserFromSessionCookie(r *http.Request, rdb *sql.DB) (*model.User, error) {
	var user *model.User
	cookieData, err := decodeSessionCookie(r)
	if err != nil {
		return user, fmt.Errorf("error getting user from session cookie: %+v", err)
	}
	sessionID, ok := cookieData["session_id"].(string)
	if !ok {
		return user, errors.New("session id not found")
	}

	if err := data.WithTx(rdb, func(tx *sql.Tx) error {
		user, err = data.ReadUserBySessionID(context.Background(), tx, sessionID)
		return err
	}); err != nil {
		return user, fmt.Errorf("err reading user by session id: %+v", err)
	}

	if !user.Expires.Valid || user.Expires.Time.Before(time.Now().UTC()) {
		return user, errors.New("session doesn't exist or has expired")
	}

	return user, nil
}

func decodeSessionCookie(r *http.Request) (map[string]any, error) {
	data := map[string]any{}
	c, err := r.Cookie(internal.SessionCookie)
	if err != nil {
		return data, err
	}
	if err := secureCookie.Decode(internal.SessionCookie, c.Value, &data); err != nil {
		return data, err
	}
	return data, nil
}

func (h *Handler) URLIDMiddleware(name string) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			idStr := c.Param(name)
			id, err := strconv.Atoi(idStr)
			if err != nil {
				return newToastError(http.StatusUnprocessableEntity, fmt.Sprintf("Invalid %s '%s'", name, idStr))
			}
			c.Set(name, id)
			return next(c)
		}
	}
}

func getURLID(c echo.Context, name string) int {
	if v, ok := c.Get(name).(int); ok {
		return v
	}
	return 0
}

func (h *Handler) RoleMiddleware(roleID internal.RoleID) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			u := getCtxUser(c)

			if u.RoleID < roleID {
				return echo.NewHTTPError(http.StatusForbidden, "Invalid permissions.")
			}

			return next(c)
		}
	}
}
