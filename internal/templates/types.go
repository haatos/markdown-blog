package templates

import "github.com/haatos/markdown-blog/internal/model"

type Page struct {
	User *model.User
	Head Head
}

type Head struct {
	Title string
	Path  string
}

func DefaultPage(user *model.User, path string) Page {
	return Page{
		User: user,
		Head: NewHead("App", path),
	}
}

func NewHead(title, path string) Head {
	return Head{
		Title: title,
		Path:  path,
	}
}

type StatusPage struct {
	Page
	Status
}

type Status struct {
	Title       string
	Description string
	Code        int
}
