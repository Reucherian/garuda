package handler

import (
	"garuda.com/m/frontend/views/layouts"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func Render(c echo.Context, status int, t templ.Component) error {
	// response header
	c.Response().Writer.WriteHeader(status)
	// response header content type
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)

	isAlpine := c.Request().Header.Get("X-Requested-With") == "AlpineJS"
	route := c.Path()

	if !isAlpine {
		return layouts.BaseLayout(route).Render(c.Request().Context(), c.Response().Writer)
	}

	return t.Render(c.Request().Context(), c.Response().Writer)
}
