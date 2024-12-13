package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/haatos/markdown-blog/internal/data"
	"github.com/haatos/markdown-blog/internal/model"
	"github.com/haatos/markdown-blog/internal/templates"
	"github.com/labstack/echo/v4"
)

type ArticlesPage struct {
	templates.Page
	Articles   []model.Article
	PageNumber int
	HasMore    bool
}

func (h *Handler) GetArticlesPage(c echo.Context) error {
	dp := getDefaultPage(c)

	var articles []model.Article
	var err error
	if dp.User.IsSuperuser() {
		articles, err = data.ReadAllArticles(c.Request().Context(), h.rdb, 6, 0, "")
		if err != nil {
			log.Println("err reading all articles", err)
		}
	} else {
		articles, err = data.ReadPublicArticles(c.Request().Context(), h.rdb, 6, 0, "")
		if err != nil {
			log.Println("err reading public articles", err)
		}
	}

	page := ArticlesPage{
		Page:     dp,
		Articles: articles,
		HasMore:  len(articles) == 6,
	}

	tn := "articles"
	if isHXRequest(c) {
		tn += "-main"
	}

	return c.Render(http.StatusOK, tn, page)
}

func (h *Handler) GetArticlesGrid(c echo.Context) error {
	dp := getDefaultPage(c)
	pageNumberStr := c.QueryParam("page")
	pageNumber, err := strconv.Atoi(pageNumberStr)
	if err != nil {
		return c.Render(http.StatusUnprocessableEntity, "toast-error", templates.NewError("Invalid page number"))
	}

	var articles []model.Article
	if dp.User.IsSuperuser() {
		articles, err = data.ReadAllArticles(c.Request().Context(), h.rdb, 6, 6*pageNumber, "")
		if err != nil {
			log.Println("err reading all articles", err)
		}
	} else {
		articles, err = data.ReadPublicArticles(c.Request().Context(), h.rdb, 6, 6*pageNumber, "")
		if err != nil {
			log.Println("err reading public articles", err)
		}
	}

	page := ArticlesPage{
		Page:       dp,
		Articles:   articles,
		PageNumber: pageNumber,
		HasMore:    len(articles) == 6,
	}

	return c.Render(http.StatusOK, "articles-grid", page)
}
