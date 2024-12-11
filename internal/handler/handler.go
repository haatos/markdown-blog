package handler

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"

	"github.com/haatos/markdown-blog/internal/templates"
	"github.com/labstack/echo/v4"
)

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
	return templates.DefaultPage(c.Request().URL.Path)
}

func generateRandomSessionID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}
