package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/georgysavva/scany/v2/sqlscan"
	"github.com/haatos/markdown-blog/internal"
	"github.com/haatos/markdown-blog/internal/model"
)

func CreateArticleComment(ctx context.Context, q sqlscan.Querier, c *model.Comment) error {
	return sqlscan.Get(ctx, q, c,
		`
        INSERT INTO comments (
			article_id,
			user_id,
			comment_id,
			content
		)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_on, updated_on
        `,
		c.ArticleID,
		c.UserID,
		c.CommentID,
		c.Content,
	)
}

func ReadArticleComments(ctx context.Context, q sqlscan.Querier, slug string, limit, offset int) ([]model.Comment, error) {
	comments := make([]model.Comment, 0, internal.CommentsPerPage)
	err := sqlscan.Select(ctx, q, &comments,
		`
        SELECT
			c.id,
			c.article_id,
			c.user_id,
			c.comment_id,
			c.content,
			c.deleted,
			c.created_on,
			c.updated_on,
			u.first_name,
			u.last_name,
			u.avatar_url
		FROM articles a
        INNER JOIN comments c ON a.id = c.article_id
		INNER JOIN users u ON c.user_id = u.id
		WHERE a.slug = $1 AND c.comment_id IS NULL
		ORDER BY c.created_on DESC
		LIMIT $2 OFFSET $3
        `,
		slug, limit, offset,
	)
	return comments, err
}

func ReadArticleCommentReplies(ctx context.Context, q sqlscan.Querier, slug string, commentID, limit, offset int) ([]model.Comment, error) {
	comments := make([]model.Comment, 0, internal.CommentsPerPage)
	err := sqlscan.Select(ctx, q, &comments,
		`
        SELECT
			c.id,
			c.article_id,
			c.user_id,
			c.comment_id,
			c.content,
			c.deleted,
			c.created_on,
			c.updated_on,
			u.first_name,
			u.last_name,
			u.avatar_url
		FROM articles a
        INNER JOIN comments c ON a.id = c.article_id
		INNER JOIN users u ON c.user_id = u.id
		WHERE a.slug = $1 AND c.comment_id = $2
		ORDER BY c.created_on DESC
		LIMIT $3 OFFSET $4
        `,
		slug, commentID, limit, offset,
	)
	return comments, err
}

func UpdateArticleComment(ctx context.Context, tx *sql.Tx, c *model.Comment) error {
	_, err := tx.ExecContext(
		ctx,
		`
        UPDATE comments
        SET content = $1,
            updated_on = $2
        WHERE id = $3
        `,
		c.Content, time.Now().UTC().Format(internal.DBTimestampLayout), c.ID,
	)
	return err
}

func DeleteArticleComment(ctx context.Context, tx *sql.Tx, commentID int, roleID internal.RoleID) error {
	var msg string
	switch roleID {
	case internal.Superuser:
		msg = "deleted by an admin"
	default:
		msg = "deleted by user"
	}
	_, err := tx.ExecContext(
		ctx,
		`
        UPDATE comments
        SET content = $1,
            updated_on = $2,
			deleted = TRUE
        WHERE id = $3
		RETURNING content, updated_on
        `,
		msg,
		time.Now().UTC().Format(internal.DBTimestampLayout),
		commentID,
	)
	return err
}

func ReadArticleCommentByID(ctx context.Context, q sqlscan.Querier, c *model.Comment) error {
	return sqlscan.Get(
		ctx, q, c,
		`
		SELECT
			c.id,
			c.article_id,
			c.user_id,
			c.content,
			c.deleted,
			c.created_on,
			c.updated_on,
			u.first_name,
			u.last_name,
			u.avatar_url
		FROM comments c
        INNER JOIN users u
		ON c.user_id = u.id
		WHERE c.id = $1
		`,
		c.ID,
	)
}

func ReadLatestArticleDiscussionComments(ctx context.Context, q sqlscan.Querier, limit int) ([]model.Comment, error) {
	comments := []model.Comment{}
	err := sqlscan.Select(
		ctx, q, &comments,
		`
		SELECT
			c.id,
			c.article_id,
			c.user_id,
			c.content,
			c.deleted,
			c.created_on,
			c.updated_on,
			u.first_name,
			u.last_name,
			u.avatar_url
        FROM comments c
		INNER JOIN users u
		ON c.user_id = u.id
		ORDER BY c.created_on DESC
		LIMIT $1
		`,
		limit,
	)
	return comments, err
}
