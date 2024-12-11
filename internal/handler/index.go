package handler

import (
	"net/http"

	"github.com/haatos/markdown-blog/internal"
	"github.com/haatos/markdown-blog/internal/templates"
	"github.com/labstack/echo/v4"
)

func (h *Handler) GetIndexPage(c echo.Context) error {
	page := getDefaultPage(c)
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
