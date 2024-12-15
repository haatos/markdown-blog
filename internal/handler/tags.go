package handler

import (
	"net/http"

	"github.com/haatos/markdown-blog/internal/data"
	"github.com/haatos/markdown-blog/internal/model"
	"github.com/haatos/markdown-blog/internal/templates"
	"github.com/labstack/echo/v4"
)

type TagsPage struct {
	templates.Page
	Tags []model.Tag
}

func (h *Handler) GetTagsPage(c echo.Context) error {
	tags, err := data.ReadTags(c.Request().Context(), h.rdb)
	if err != nil {
		c.Logger().Error("err reading tags ", err)
	}

	page := TagsPage{
		Page: getDefaultPage(c),
		Tags: tags,
	}

	tn := "tags"
	if isHXRequest(c) {
		tn += "-main"
	}
	return c.Render(http.StatusOK, tn, page)
}
