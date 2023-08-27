package main

import (
	"fmt"
	"html/template"
	"htmx/handlers"
	"io"

	"github.com/labstack/echo/v4"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tmpl, err := template.New("").ParseFiles(
		"static/nav.html",
		"static/users.html",
		"static/logged-in.html",
		"static/register.html",
		"static/about.html",
		"static/index.html",
		"static/base.html",
	)
	if err != nil {
		fmt.Println("Error", err)
		return err
	}
	return tmpl.ExecuteTemplate(w, "base", data)
}

func main() {
	e := echo.New()
	e.Renderer = &Template{templates: template.Must(template.ParseGlob("static/*.html"))}
	e.Static("/static", "static")

	handlers.SeedUsers()
	// e.File("/", "static/index.html")

	e.GET("/", handlers.Index)
	e.GET("/about", handlers.About)
	e.GET("/register", handlers.RenderRegister)
	e.GET("/login", handlers.RenderLogin)
	e.GET("/users", handlers.AllUsers)
	e.POST("/register", handlers.RegisterUser)
	e.POST("/login", handlers.Login)

	e.Logger.Fatal(e.Start(":8080"))
}
