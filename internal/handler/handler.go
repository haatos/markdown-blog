package handler

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"net/http"

	"github.com/haatos/markdown-blog/internal/model"
	"github.com/haatos/markdown-blog/internal/templates"
	"github.com/labstack/echo/v4"
)

func getCtxUser(c echo.Context) *model.User {
	u, ok := c.Get("user").(*model.User)
	if ok {
		return u
	}
	return nil
}

func NewHandler(rdb *sql.DB, rwdb *sql.DB) *Handler {
	return &Handler{
		rdb:  rdb,
		rwdb: rwdb,
	}
}

type Handler struct {
	rdb  *sql.DB
	rwdb *sql.DB
}

func templateName(c echo.Context, name string) string {
	if isHXRequest(c) {
		name += "-main"
	}
	return name
}

func getDefaultPage(c echo.Context) templates.Page {
	u := getCtxUser(c)
	return templates.DefaultPage(u, c.Request().URL.Path)
}

func generateRandomSessionID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

func (h *Handler) renderStatusPage(c echo.Context, status int) error {
	dp := getDefaultPage(c)
	tn := "status"
	if isHXRequest(c) {
		tn += "-main"
	}
	switch status {
	case http.StatusNotFound:
		return c.Render(
			http.StatusInternalServerError, tn,
			templates.StatusPage{
				Page: dp,
				Status: templates.Status{
					Code:        http.StatusNotFound,
					Title:       "Not Found",
					Description: "Seems like there's nothing here...",
				},
			})
	case http.StatusForbidden:
		return c.Render(
			http.StatusInternalServerError, tn,
			templates.StatusPage{
				Page: dp,
				Status: templates.Status{
					Code:        http.StatusForbidden,
					Title:       "Forbidden",
					Description: "Invalid permissions to view this content.",
				},
			})
	default:
		return c.Render(
			http.StatusInternalServerError, tn,
			templates.StatusPage{
				Page: dp,
				Status: templates.Status{
					Code:        http.StatusInternalServerError,
					Title:       "Internal Server Error",
					Description: "Something went terribly wrong.",
				},
			})
	}
}

func (h *Handler) renderToastInfo(c echo.Context, message string) error {
	hxRetarget(c, "body")
	hxReswap(c, "beforeend")
	return c.Render(http.StatusOK, "toast-info", message)
}
