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
	// "logged-in.html": {"templates/shared/logged-in.html"},
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tmp := template.Must(template.ParseGlob("templates/shared/*.html"))

	f := tmps[name]
	if len(f) > 0 {
		if err := template.Must(tmp.ParseFiles(f...)).ExecuteTemplate(w, "base", data); err != nil {
			fmt.Println("Error", err)
			return err
		}
	} else {
		if err := tmp.ExecuteTemplate(w, name, data); err != nil {
			fmt.Println("Error", err)
			return err
		}
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
