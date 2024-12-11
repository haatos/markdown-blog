package templates

type Page struct {
	Head Head
}

type Head struct {
	Title string
	Path  string
}

func DefaultPage(path string) Page {
	return Page{
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
