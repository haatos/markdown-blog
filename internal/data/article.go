package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/georgysavva/scany/v2/sqlscan"
	"github.com/haatos/markdown-blog/internal/model"
)

func CreateArticle(ctx context.Context, q sqlscan.Querier, article *model.Article) error {
	return sqlscan.Get(ctx, q, article,
		`
		INSERT INTO articles (
			user_id,
			title,
			slug,
			description,
			content
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
		`,
		article.UserID,
		article.Title,
		article.Slug,
		article.Description,
		article.Content,
	)
}

func ReadArticleByID(ctx context.Context, q sqlscan.Querier, article *model.Article) error {
	return sqlscan.Get(ctx, q, article,
		`
		SELECT
			user_id,
			title,
			slug,
			description,
			content,
			published_on
		FROM articles
		WHERE id = $1
		`,
		article.ID,
	)
}

func DeleteArticle(ctx context.Context, db *sql.DB, id int) error {
	_, err := db.Exec(`DELETE FROM articles WHERE id = $1`, id)
	return err
}

func ReadArticleBySlug(ctx context.Context, q sqlscan.Querier, article *model.Article) error {
	return sqlscan.Get(ctx, q, article,
		`
		SELECT
			a.id,
			a.title,
			a.slug,
			a.description,
			a.content,
			a.published_on,
			u.first_name,
			u.last_name
		FROM articles a
		INNER JOIN users u
		ON a.user_id = u.id
		WHERE a.slug = $1
		`,
		article.Slug,
	)
}

func ReadAllArticles(ctx context.Context, q sqlscan.Querier, limit, offset int, filter string) ([]model.Article, error) {
	as := make([]model.Article, 0, limit)
	err := sqlscan.Select(
		ctx, q, &as,
		`
		SELECT
			id,
			user_id,
			title,
			slug,
			description,
			content,
			published_on
		FROM articles
		WHERE LOWER(title) LIKE '%'||$1||'%'
		ORDER BY articles.published_on DESC
		LIMIT $2 OFFSET $3
		`,
		filter, limit, offset,
	)
	return as, err
}

func ReadAllImages(ctx context.Context, q sqlscan.Querier) ([]model.Image, error) {
	images := make([]model.Image, 0)
	err := sqlscan.Select(ctx, q, &images, "select id, name, image_key from images")
	return images, err
}

func ReadPublicArticles(ctx context.Context, q sqlscan.Querier, limit, offset int, filter string) ([]model.Article, error) {
	as := make([]model.Article, 0, limit)
	err := sqlscan.Select(
		ctx, q, &as,
		`
		SELECT
			id,
			user_id,
			title,
			slug,
			description,
			published_on
		FROM articles
		WHERE published_on IS NOT NULL AND LOWER(title) LIKE '%'||$1||'%'
		ORDER BY articles.published_on DESC
		LIMIT $2 OFFSET $3
		`,
		filter, limit, offset,
	)
	return as, err
}

func ReadLatestArticles(ctx context.Context, q sqlscan.Querier, limit int) ([]model.Article, error) {
	as := make([]model.Article, 0, limit)
	err := sqlscan.Select(
		ctx, q, &as,
		`
		SELECT
			id,
			user_id,
			title,
			slug,
			description,
			published_on
		FROM articles
		WHERE published_on IS NOT NULL
		ORDER BY published_on DESC
		LIMIT $1
		`,
		limit,
	)
	return as, err
}

func handleArticleRows(as []model.Article, rows *sql.Rows) ([]model.Article, error) {
	var articleID int
	var a model.Article
	var tags []model.Tag
	for rows.Next() {
		var aID, aUserID int
		var aTitle, aSlug, aDescription, aImageKey string
		var aUpdatedOn time.Time
		var aPublishedOn *time.Time
		var tagID *int
		var tagName *string

		if err := rows.Scan(
			&aID,
			&aUserID,
			&aTitle,
			&aSlug,
			&aDescription,
			&aImageKey,
			&aPublishedOn,
			&aUpdatedOn,
			&tagID,
			&tagName,
		); err != nil {
			return as, err
		}

		if articleID == 0 {
			articleID = aID
			a = model.Article{
				ID:          aID,
				UserID:      aUserID,
				Title:       aTitle,
				Slug:        aSlug,
				Description: aDescription,
				ImageKey:    aImageKey,
				PublishedOn: aPublishedOn,
			}
			tags = make([]model.Tag, 0, 3)
		} else if aID != articleID {
			articleID = aID
			a.Tags = tags
			as = append(as, a)
			tags = make([]model.Tag, 0, 3)
			a = model.Article{
				ID:          aID,
				UserID:      aUserID,
				Title:       aTitle,
				Slug:        aSlug,
				Description: aDescription,
				ImageKey:    aImageKey,
				PublishedOn: aPublishedOn,
			}
		}
		if tagID != nil && tagName != nil {
			tags = append(tags, model.Tag{ID: *tagID, Name: *tagName})
		}
	}
	if a.ID != 0 {
		a.Tags = tags
		as = append(as, a)
	}
	return as, nil
}

func ReadArticleTags(ctx context.Context, q sqlscan.Querier, a *model.Article) ([]model.Tag, error) {
	tags := make([]model.Tag, 0, 3)
	err := sqlscan.Select(
		ctx, q, &tags,
		`
		SELECT t.id, t.name
		FROM tags t
		INNER JOIN articles_tags at ON t.id = at.tag_id
		WHERE at.article_id = $1
		`,
		a.ID,
	)
	return tags, err
}

func DeleteArticleTags(ctx context.Context, tx *sql.Tx, articleID int) error {
	_, err := tx.ExecContext(
		ctx,
		`
		DELETE FROM articles_tags WHERE article_id = $1
		`,
		articleID,
	)
	return err
}

func ReadRelatedArticlesByID(ctx context.Context, q sqlscan.Querier, id int, limit int) ([]model.Article, error) {
	as := make([]model.Article, 0, limit)
	rows, err := q.QueryContext(
		ctx,
		`
		SELECT
			a.id,
			a.user_id,
			a.title,
			a.slug,
			a.description,
			a.published_on,
			t.id,
			t.name
		FROM articles a
		LEFT JOIN articles_tags at ON a.id = at.article_id
		LEFT JOIN tags t ON at.tag_id = t.id
		WHERE t.id IN (
			SELECT tag_id
			FROM articles_tags
			WHERE article_id = $1
		)
		AND a.id != $1
		GROUP BY a.id, a.title, a.content
		ORDER BY COUNT(*) DESC, published_on DESC
		LIMIT $2
		`,
		id,
		limit,
	)
	if err != nil {
		return as, err
	}
	defer rows.Close()

	return handleArticleRows(as, rows)
}
