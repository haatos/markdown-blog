package model

import (
	"fmt"

	"github.com/haatos/markdown-blog/internal"
)

type Image struct {
	ID       int
	Name     string
	ImageKey string
}

func (i *Image) ImageURL() string {
	if i.ImageKey == "" {
		return ""
	}
	return fmt.Sprintf("https://%s%s", internal.Settings.CloudfrontDomain, i.ImageKey)
}
