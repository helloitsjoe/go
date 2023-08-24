package main

import (
	"htmx/handlers"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.Static("/static", "static")
	e.File("/", "static/index.html")

	e.POST("/register", handlers.RegisterUser)
	e.POST("/login", handlers.Login)

	e.Logger.Fatal(e.Start(":8080"))
}
