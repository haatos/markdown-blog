package handler

import (
	"net/http"
	"strings"

	"github.com/haatos/markdown-blog/internal/templates"
	"github.com/labstack/echo/v4"
)

type ErrorPage struct {
	templates.Page
	Code        int
	Title       string
	Description string
}

func ErrorHandler(err error, c echo.Context) {
	c.Logger().Errorf("Handler error: %+v\n", err)
	switch e := err.(type) {
	case ToastError:
		c.Render(e.Status, "toast-error", e.Messages)
	case *echo.HTTPError:
		switch e.Code {
		case http.StatusNotFound:
			if isHXRequest(c) {
				c.Render(
					http.StatusNotFound, "status-main",
					ErrorPage{
						Page:        getDefaultPage(c),
						Code:        http.StatusNotFound,
						Title:       "Not Found",
						Description: "There seems to be nothing here.",
					})
			} else {
				c.Render(
					http.StatusNotFound, "status",
					ErrorPage{
						Page:        getDefaultPage(c),
						Code:        http.StatusNotFound,
						Title:       "Not Found",
						Description: "There seems to be nothing here.",
					})
			}
			return
		case http.StatusInternalServerError:
			if isHXRequest(c) {
				c.Render(
					http.StatusNotFound, "status-main",
					ErrorPage{
						Page:        getDefaultPage(c),
						Code:        http.StatusInternalServerError,
						Title:       "Internal Server Error",
						Description: "Something went terribly wrong. Please try again later.",
					})
			} else {
				c.Render(
					http.StatusNotFound, "status",
					ErrorPage{
						Page:        getDefaultPage(c),
						Code:        http.StatusInternalServerError,
						Title:       "Internal Server Error",
						Description: "Something went terribly wrong. Please try again later.",
					})
			}
			return
		case http.StatusForbidden:
			if isHXRequest(c) {
				c.Render(
					http.StatusNotFound, "status-main",
					ErrorPage{
						Page:        getDefaultPage(c),
						Code:        http.StatusForbidden,
						Title:       "Forbidden",
						Description: "Invalid permissions to view this page.",
					})
			} else {
				c.Render(
					http.StatusNotFound, "status",
					ErrorPage{
						Page:        getDefaultPage(c),
						Code:        http.StatusForbidden,
						Title:       "Forbidden",
						Description: "Invalid permissions to view this page.",
					})
			}
			return
		}
	}
}

func newToastError(status int, messages ...string) ToastError {
	return ToastError{
		Status:   status,
		Messages: messages,
	}
}

type ToastError struct {
	Status   int
	Messages []string
}

func (te ToastError) Error() string {
	return strings.Join(te.Messages, ", ")
}

func logError(c echo.Context, message string, err error) {
	c.Logger().Printf("%s: %+v\n", message, err)
}
