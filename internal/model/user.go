package model

import (
	"database/sql"

	"github.com/haatos/markdown-blog/internal"
)

type User struct {
	ID        int
	RoleID    internal.RoleID
	Email     string
	FirstName string
	LastName  string
	AvatarURL string

	// session
	Expires sql.NullTime
}

func (u *User) IsOwner() bool {
	return u != nil && u.Email == internal.Settings.SuperuserEmail
}
