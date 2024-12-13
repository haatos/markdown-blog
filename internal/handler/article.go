package handler

import (
	"fmt"
	"net/http"

	"github.com/haatos/markdown-blog/internal/data"
	"github.com/haatos/markdown-blog/internal/model"
	"github.com/haatos/markdown-blog/internal/templates"
	"github.com/labstack/echo/v4"
)

type ArticlePage struct {
	templates.Page
	Article model.Article
}

func (h *Handler) GetArticlePage(c echo.Context) error {
	dp := getDefaultPage(c)

	article := model.Article{Slug: c.Param("slug")}
	err := data.ReadArticleBySlug(c.Request().Context(), h.rdb, &article)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("err reading article by slug: %+v\n", err))
	}

	article.Content = getHTMLFromMarkdown([]byte(article.Content))

	page := ArticlePage{
		Page:    dp,
		Article: article,
	}

	tn := "article"
	if isHXRequest(c) {
		tn += "-main"
	}
	return c.Render(http.StatusOK, tn, page)
}
