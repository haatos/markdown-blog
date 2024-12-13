package templates

import (
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"path/filepath"
	"strings"
)

func NewRenderer() *Template {
	paths := make([]string, 0, 10)
	filepath.Walk("./internal/templates", func(path string, info fs.FileInfo, err error) error {
		if strings.HasSuffix(path, ".html") {
			paths = append(paths, path)
		}
		return nil
	})
	t := template.Must(template.New("").Funcs(template.FuncMap{
		"increment": func(v, by int) int { return v + by },
		"hasPrefix": strings.HasPrefix,
		"rawHTML":   func(s string) template.HTML { return template.HTML(s) },
		"N": func(n int) []struct{} {
			return make([]struct{}, n)
		},
		"withAttrs": func(pairs ...any) (map[string]any, error) {
			if len(pairs)%2 != 0 {
				return nil, errors.New("input argument count must be even")
			}
			attrs := make(map[string]any, len(pairs)/2)
			for i := 0; i < len(pairs); i += 2 {
				k := fmt.Sprintf("%v", pairs[i])
				v := pairs[i+1]

				if k == "attrs" {
					attrs, ok := v.(map[string]any)
					if ok {
						for kk, vv := range attrs {
							attrs[kk] = vv
						}
						continue
					}
				}

				attrs[k] = v
			}

			return attrs, nil
		},
	}).ParseFiles(paths...))
	return &Template{templates: t}
}

type Template struct {
	templates *template.Template
}

type Error struct {
	Messages []string
}

func NewError(messages ...string) Error {
	return Error{
		Messages: messages,
	}
}
