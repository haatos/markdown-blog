package handler

import (
	"net/http"
	"strings"

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

func (h *Handler) PostTags(c echo.Context) error {
	name := c.FormValue("name")
	name = strings.TrimSpace(name)

	if name == "" {
		return newToastError(http.StatusUnprocessableEntity, "Tag name cannot be empty")
	}

	tag := model.Tag{Name: name}
	if err := data.CreateTag(c.Request().Context(), h.rwdb, &tag); err != nil {
		c.Logger().Error("err creating tag: ", err)
		return newToastError(http.StatusInternalServerError, "Unable to create tag")
	}

	return c.Render(http.StatusOK, "tag-div", tag)
}

func (h *Handler) DeleteTag(c echo.Context) error {
	id := getURLID(c, "tagID")

	if err := data.DeleteTag(c.Request().Context(), h.rwdb, id); err != nil {
		c.Logger().Error("err deleting tag: ", err)
		return newToastError(http.StatusInternalServerError, "Unable to delete tag")
	}

	return c.String(http.StatusOK, "")
}
