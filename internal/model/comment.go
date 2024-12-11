package model

import (
	"database/sql"
	"fmt"
	"time"
)

type Comment struct {
	ID        int
	ArticleID int
	UserID    int
	CommentID sql.NullInt64
	Content   string
	Deleted   bool
	CreatedOn time.Time
	UpdatedOn time.Time

	FirstName string
	LastName  string
	AvatarURL string
}

func (ac *Comment) CreatedAgo() string {
	return durationToString(time.Since(ac.CreatedOn).Round(time.Minute))
}

func (ac *Comment) FmtCreatedOn() string {
	return ac.CreatedOn.Format("Jan 02 2006")
}

func (ac *Comment) UpdatedAgo() string {
	return durationToString(time.Since(ac.UpdatedOn).Round(time.Minute))
}

func durationToString(d time.Duration) string {
	s := d.Milliseconds() / 1000

	result := ""
	years := s / (365 * 24 * 60 * 60)
	if years > 0 {
		result += fmt.Sprintf("%dy", years)
		s -= 365 * 24 * 60 * 60 * years
	}

	days := s / (24 * 60 * 60)
	if days > 0 {
		result += fmt.Sprintf(" %dd", days)
		return result
	}

	hours := s / (60 * 60)
	if hours > 0 {
		result += fmt.Sprintf(" %dh", hours)
		return result
	}

	minutes := s / 60
	if minutes > 0 {
		result += fmt.Sprintf(" %dm", minutes)
	}

	return result
}
