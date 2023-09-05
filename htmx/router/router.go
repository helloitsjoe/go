package router

import (
	"htmx/templates"

	"github.com/labstack/echo/v4"
)

func New(rootDir string) *echo.Echo {
	e := echo.New()
	e.Renderer = &templates.Template{RootDir: rootDir}

	e.Static("/static", "static")

	return e
}
