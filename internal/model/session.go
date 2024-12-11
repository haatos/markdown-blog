package model

import "time"

type Session struct {
	ID      string
	UserID  int
	Expires time.Time
}
