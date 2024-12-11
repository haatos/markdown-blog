package handler

import "github.com/labstack/echo/v4"

func hxRedirect(c echo.Context, path string) {
	c.Response().Writer.Header().Set("HX-Redirect", path)
}

func hxRetarget(c echo.Context, target string) {
	c.Response().Writer.Header().Set("HX-Retarget", target)
}

func hxReswap(c echo.Context, swap string) {
	c.Response().Writer.Header().Set("HX-Reswap", swap)
}

func hxPushURL(c echo.Context, url string) {
	c.Response().Writer.Header().Set("HX-Push-Url", url)
}

func isHXRequest(c echo.Context) bool {
	return c.Request().Header.Get("hx-request") != ""
}
