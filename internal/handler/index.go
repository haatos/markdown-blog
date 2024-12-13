package handler

import (
	"log"
	"net/http"

	"github.com/haatos/markdown-blog/internal"
	"github.com/haatos/markdown-blog/internal/data"
	"github.com/haatos/markdown-blog/internal/model"
	"github.com/haatos/markdown-blog/internal/templates"
	"github.com/labstack/echo/v4"
)

type IndexPage struct {
	templates.Page
	RecentArticles []model.Article
}

func (h *Handler) GetIndexPage(c echo.Context) error {
	dp := getDefaultPage(c)
	articles, err := data.ReadLatestArticles(c.Request().Context(), h.rdb, 3)
	if err != nil {
		log.Println("err reading latest articles", err)
	}
	page := IndexPage{
		Page:           dp,
		RecentArticles: articles,
	}
	return c.Render(http.StatusOK, templateName(c, "index"), page)
}

type LegalPage struct {
	templates.Page
	Domain       string
	ContactEmail string
}

func (h *Handler) GetPrivacyPage(c echo.Context) error {
	page := LegalPage{
		Page:         getDefaultPage(c),
		Domain:       internal.Settings.Domain,
		ContactEmail: internal.Settings.ContactEmail,
	}
	return c.Render(http.StatusOK, templateName(c, "privacy"), page)
}

func (h *Handler) GetTermsOfServicePage(c echo.Context) error {
	page := LegalPage{
		Page:         getDefaultPage(c),
		Domain:       internal.Settings.Domain,
		ContactEmail: internal.Settings.ContactEmail,
	}
	return c.Render(http.StatusOK, templateName(c, "terms-of-service"), page)
}
