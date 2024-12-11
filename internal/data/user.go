package data

import (
	"context"
	"database/sql"
	"strings"

	"github.com/georgysavva/scany/v2/sqlscan"
	"github.com/haatos/markdown-blog/internal"
	"github.com/haatos/markdown-blog/internal/model"
)

func CreateUser(ctx context.Context, tx sqlscan.Querier, user *model.User) error {
	if user.Email == internal.Settings.SuperuserEmail {
		user.RoleID = internal.Superuser
	}

	return sqlscan.Get(
		ctx, tx, user,
		`
		INSERT INTO users (
			role_id,
			first_name,
			last_name,
			email,
			avatar_url
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
		`,
		user.RoleID,
		user.FirstName,
		user.LastName,
		user.Email,
		user.AvatarURL,
	)
}

func ReadUserByID(ctx context.Context, qr sqlscan.Querier, user *model.User) error {
	return sqlscan.Get(
		ctx, qr, user,
		`
		SELECT
			role_id,
			first_name,
			last_name,
			email,
			avatar_url
		FROM users WHERE id = $1
		`,
		user.ID,
	)
}

func ReadUserByEmail(ctx context.Context, q sqlscan.Querier, user *model.User) error {
	return sqlscan.Get(ctx, q, user,
		`
		SELECT
			id,
			role_id,
			first_name,
			last_name,
			avatar_url
		FROM users WHERE email = $1
		`,
		user.Email,
	)
}

func ReadUsers(ctx context.Context, q sqlscan.Querier, limit, offset int, filter string) ([]model.User, error) {
	users := make([]model.User, 0, limit)
	filter = strings.ToLower(filter)

	err := sqlscan.Select(ctx, q, &users,
		`
			SELECT
				u.id,
				u.role_id,
				u.email,
				u.first_name,
				u.last_name,
				u.avatar_url
			FROM users u
			WHERE LOWER(email) LIKE '%'||$1||'%'
			OR LOWER(first_name) LIKE '%'||$1||'%'
			OR LOWER(last_name) LIKE '%'||$1||'%'
			ORDER BY last_name ASC
			LIMIT $2 OFFSET $3
		`,
		filter, limit, offset,
	)
	return users, err
}

func ReadUserBySessionID(ctx context.Context, q sqlscan.Querier, sessionID string) (*model.User, error) {
	user := &model.User{}
	err := sqlscan.Get(
		ctx, q, user,
		`
		SELECT
			u.id,
			u.role_id,
			u.email,
			u.first_name,
			u.last_name,
			u.avatar_url,
			s.expires
		FROM users u
		LEFT JOIN sessions s
		ON u.id = s.user_id
		WHERE s.id = $1
		ORDER BY s.expires DESC
		LIMIT 1
		`,
		sessionID,
	)
	return user, err
}

func DeleteUser(ctx context.Context, tx *sql.Tx, userID int) error {
	_, err := tx.ExecContext(
		ctx,
		`
		DELETE FROM users WHERE id = $1
		`,
		userID,
	)
	return err
}
