package main

import (
	"htmx/handlers"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.Static("/static", "static")
	e.GET("/", handlers.Index)
	e.POST("/api", handlers.Update)
	e.Logger.Fatal(e.Start(":8080"))
}
