package main

import (
	"fmt"
	"html/template"
	"htmx/handlers"
	"io"

	"github.com/labstack/echo/v4"
)

type Template struct{}

var tmps = map[string][]string{
	"index.html": {"templates/index.html"},
	"about.html": {"templates/about.html"},
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	f := tmps[name]
	if len(f) == 0 {
		f = []string{}
	}

	// TODO: if no files, don't call parsefiles
	if err := template.Must(template.Must(template.ParseGlob("templates/shared/*.html")).ParseFiles(f...)).ExecuteTemplate(w, "base", data); err != nil {
		fmt.Println("Error", err)
		return err
	}
	return nil
}

func main() {
	e := echo.New()
	e.Renderer = &Template{}

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
