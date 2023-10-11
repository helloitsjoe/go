package handlers

import (
	"github.com/labstack/echo/v4"

	"htmx/db"
	"htmx/middleware"
	"htmx/user"
)

// Typecheck in editor is not working correctly
func Register(e *echo.Echo) {
	d := db.CreateDB()
	user.SeedUsers(d)
	h := NewHandlers(d)

	e.GET("/", middleware.Auth(h.Index))
	e.GET("/about", h.About)
	e.GET("/register", h.RenderRegister)
	e.GET("/login", middleware.Auth(h.RenderLogin))
	e.POST("/logout", h.Logout)
	e.GET("/users", h.AllUsers)
	e.POST("/register", h.RegisterUser)
	e.POST("/login", h.Login)
}
