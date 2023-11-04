package handlers

import (
	"github.com/labstack/echo/v4"

	"htmx/db"
	"htmx/middleware"
)

func Register(e *echo.Echo, d db.DB) {
	h := NewHandlers(d)

	e.GET("/", middleware.Auth(h.Index))
	e.GET("/register", middleware.Auth(h.Index))
	e.GET("/login", middleware.Auth(h.Index))
	e.GET("/followers", middleware.Auth(h.MyFollowers))
	e.GET("/about", h.About)
	e.POST("/logout", h.Logout)
	e.GET("/users", h.AllUsers)
	e.POST("/register", h.RegisterUser)
	e.POST("/login", h.Login)
}
