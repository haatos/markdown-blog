package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/haatos/markdown-blog/internal"
	"github.com/haatos/markdown-blog/internal/data"
	"github.com/haatos/markdown-blog/internal/model"
	"github.com/haatos/markdown-blog/internal/templates"
	"github.com/labstack/echo/v4"
)

type ArticleEditorPage struct {
	templates.Page
	Article           model.Article
	ArticleTagsSelect ArticleTagsSelect
}

func (h *Handler) GetArticleEditorPage(c echo.Context) error {
	id := getURLID(c, "articleID")
	article := model.Article{ID: id}

	if err := data.ReadArticleByID(c.Request().Context(), h.rwdb, &article); err != nil {
		c.Logger().Error("err reading article by ID", err)
	}

	if err := data.ReadArticleTags(c.Request().Context(), h.rdb, &article); err != nil {
		c.Logger().Error("err reading article tags", err)
	}

	tags, err := data.ReadTags(c.Request().Context(), h.rdb)
	if err != nil {
		c.Logger().Error("err reading tags", err)
	}

	page := ArticleEditorPage{
		Page:    getDefaultPage(c),
		Article: article,
		ArticleTagsSelect: ArticleTagsSelect{
			ArticleTags: article.Tags,
			Tags:        tags,
		},
	}

	return c.Render(http.StatusOK, "article-editor", page)
}

func (h *Handler) PatchArticleTitle(c echo.Context) error {
	id := getURLID(c, "articleID")
	title := c.FormValue("title")

	_, err := h.rwdb.ExecContext(c.Request().Context(), data.UpdateIDQuery("articles", "title"), title, id)
	if err != nil {
		c.Logger().Error("err updating article title", err)
		return newToastError(http.StatusInternalServerError, "Unable to update article title.")
	}

	return h.renderToastInfo(c, "Title updated")
}

func (h *Handler) PatchArticleDescription(c echo.Context) error {
	id := getURLID(c, "articleID")
	description := c.FormValue("description")

	_, err := h.rwdb.ExecContext(c.Request().Context(), data.UpdateIDQuery("articles", "description"), description, id)
	if err != nil {
		c.Logger().Error("err updating article description", err)
		return newToastError(http.StatusInternalServerError, "Unable to update article description.")
	}

	return h.renderToastInfo(c, "Description updated")
}

func (h *Handler) PatchArticleContent(c echo.Context) error {
	id := getURLID(c, "articleID")
	content := c.FormValue("content")

	_, err := h.rwdb.ExecContext(c.Request().Context(), data.UpdateIDQuery("articles", "content"), content, id)
	if err != nil {
		c.Logger().Error("err updating article content", err)
		return newToastError(http.StatusInternalServerError, "Unable tot update article content.")
	}

	return h.renderToastInfo(c, "Content updated")
}

func (h *Handler) PatchArticleVisibility(c echo.Context) error {
	id := getURLID(c, "articleID")
	public := c.FormValue("visibility") == "on"

	var err error
	if public {
		_, err = h.rwdb.ExecContext(
			c.Request().Context(),
			data.UpdateIDQuery("articles", "published_on"),
			time.Now().UTC().Format(internal.DBTimestampLayout), id,
		)
	} else {
		_, err = h.rwdb.ExecContext(
			c.Request().Context(),
			"update articles set published_on = null where id = $1",
			id,
		)
	}
	if err != nil {
		c.Logger().Error("err updating article visibility", err)
		return newToastError(http.StatusInternalServerError, "Unable to update article visibility")
	}
	return c.Render(http.StatusOK, "article-editor-public-label", public)
}

func (h *Handler) DeleteArticle(c echo.Context) error {
	id := getURLID(c, "articleID")

	err := data.DeleteArticle(c.Request().Context(), h.rwdb, id)
	if err != nil {
		c.Logger().Error("err deleting article", err)
		return newToastError(http.StatusInternalServerError, "Unable to delete article.")
	}

	hxRedirect(c, "/articles")
	return nil
}

type ArticleTagsSelect struct {
	ArticleTags []model.Tag
	Tags        []model.Tag
}

func (h *Handler) PutArticleTags(c echo.Context) error {
	values, err := c.FormParams()
	if err != nil {
		return newToastError(http.StatusUnprocessableEntity, "Unable to parse tags")
	}

	tagIDs := values["tags"]

	id := getURLID(c, "articleID")
	a := model.Article{ID: id}
	var tags []model.Tag
	if err := data.WithTx(h.rwdb, func(tx *sql.Tx) error {
		if err := data.DeleteArticleTags(c.Request().Context(), tx, id); err != nil {
			c.Logger().Error("err deleting article tags ", err)
			return err
		}
		stmt, err := tx.Prepare(`INSERT INTO articles_tags (article_id, tag_id) VALUES ($1, $2)`)
		if err != nil {
			return err
		}
		for _, t := range tagIDs {
			fmt.Printf("inserting articles_tags (%d, %s)\n", id, t)
			if _, err := stmt.ExecContext(c.Request().Context(), id, t); err != nil {
				c.Logger().Error("err inserting articles_tags ", err)
				return err
			}
		}
		return err
	}); err != nil {
		return newToastError(http.StatusInternalServerError, "Unable to update article tags")
	}

	if err := data.ReadArticleByID(c.Request().Context(), h.rdb, &a); err != nil {
		c.Logger().Error("err reading article by id ", err)
		return err
	}

	err = data.ReadArticleTags(c.Request().Context(), h.rdb, &a)

	tags, err = data.ReadTags(c.Request().Context(), h.rdb)
	if err != nil {
		c.Logger().Error("err reading tags ", err)
	}

	return c.Render(http.StatusOK, "article-tags-select", ArticleTagsSelect{ArticleTags: a.Tags, Tags: tags})
}
