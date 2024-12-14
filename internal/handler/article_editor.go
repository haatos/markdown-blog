package handler

import (
	"net/http"

	"github.com/haatos/markdown-blog/internal/data"
	"github.com/haatos/markdown-blog/internal/model"
	"github.com/haatos/markdown-blog/internal/templates"
	"github.com/labstack/echo/v4"
)

type ArticleEditorPage struct {
	templates.Page
	Article model.Article
}

func (h *Handler) GetArticleEditorPage(c echo.Context) error {
	id := getURLID(c, "articleID")
	article := model.Article{ID: id}

	if err := data.ReadArticleByID(c.Request().Context(), h.rwdb, &article); err != nil {
		c.Logger().Error("err reading article by ID", err)
	}

	page := ArticleEditorPage{
		Page:    getDefaultPage(c),
		Article: article,
	}

	tn := "article-editor"
	if isHXRequest(c) {
		tn += "-main"
	}
	return c.Render(http.StatusOK, tn, page)
}
