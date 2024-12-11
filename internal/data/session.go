package data

import (
	"context"
	"database/sql"

	"github.com/haatos/markdown-blog/internal"
	"github.com/haatos/markdown-blog/internal/model"
)

func CreateSession(ctx context.Context, tx *sql.DB, sess *model.Session) error {
	_, err := tx.ExecContext(
		ctx,
		`
		INSERT INTO sessions (id, user_id, expires)
		VALUES ($1, $2, $3)
		`,
		sess.ID, sess.UserID, sess.Expires.Format(internal.DBTimestampLayout),
	)
	return err
}

func DeleteSessionsByUserID(ctx context.Context, tx *sql.Tx, userID int) error {
	_, err := tx.ExecContext(
		ctx,
		`
		DELETE FROM sessions WHERE user_id = $1
		`,
		userID,
	)
	return err
}
