package handler

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"

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
