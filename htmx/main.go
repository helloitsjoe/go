package main

import (
	"html/template"
	"htmx/handlers"
	"io"

	"github.com/labstack/echo/v4"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	e := echo.New()
	e.Renderer = &Template{templates: template.Must(template.ParseGlob("static/*.html"))}
	e.Static("/static", "static")
	// e.File("/", "static/index.html")

	e.GET("/", handlers.Index)
	e.POST("/register", handlers.RegisterUser)
	e.POST("/login", handlers.Login)

	e.Logger.Fatal(e.Start(":8080"))
}
