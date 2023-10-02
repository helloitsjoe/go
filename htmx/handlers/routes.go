package handlers

import (
	"github.com/labstack/echo/v4"

	"htmx/middleware"
	"htmx/user"
)

// Typecheck in editor is not working correctly
func Register(e *echo.Echo) {
	user.SeedUsers()

	e.GET("/", middleware.Auth(Index))
	e.GET("/about", About)
	e.GET("/register", RenderRegister)
	e.GET("/login", RenderLogin)
	e.GET("/users", AllUsers)
	e.POST("/register", RegisterUser)
	e.POST("/login", Login)
}
