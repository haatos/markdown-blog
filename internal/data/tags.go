package data

import (
	"context"
	"database/sql"

	"github.com/georgysavva/scany/v2/sqlscan"
	"github.com/haatos/markdown-blog/internal/model"
)

func CreateTag(ctx context.Context, q sqlscan.Querier, t *model.Tag) error {
	return sqlscan.Get(
		ctx, q, t,
		`
		INSERT INTO tags (name) VALUES ($1) RETURNING id
		`,
		t.Name,
	)
}

func ReadTags(ctx context.Context, q sqlscan.Querier) ([]model.Tag, error) {
	var tags []model.Tag
	err := sqlscan.Select(
		ctx, q, &tags,
		`
		SELECT id, name FROM tags
		`,
	)
	return tags, err
}

func DeleteTag(ctx context.Context, db *sql.DB, id int) error {
	_, err := db.ExecContext(ctx, "delete from tags where id = $1", id)
	return err
}
