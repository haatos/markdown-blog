package model

import (
	"fmt"
	"time"

	"github.com/haatos/markdown-blog/internal"
)

type Article struct {
	ID          int
	UserID      int
	Title       string
	Slug        string
	Description string
	Content     string
	ImageKey    string
	Public      bool
	PublishedOn time.Time
	CreatedOn   time.Time
	UpdatedOn   time.Time

	FirstName string
	LastName  string
	Tags      []Tag
}

func (a *Article) ImageURL() string {
	if a.ImageKey == "" {
		return ""
	}
	return fmt.Sprintf("https://%s%s", internal.Settings.CloudfrontDomain, a.ImageKey)
}
