package handlers

import (
	"github.com/labstack/echo/v4"

	"htmx/db"
	"htmx/middleware"
)

func Register(e *echo.Echo, d db.DB) {
	h := NewHandlers(d)

	e.GET("/", middleware.Auth(h.Index, d))
	e.GET("/register", middleware.Auth(h.Index, d))
	e.GET("/login", middleware.Auth(h.Index, d))
	e.GET("/followers", middleware.Auth(h.RenderFollowers, d))
	e.GET("/following", middleware.Auth(h.RenderFollowing, d))
	e.GET("/user/:username", middleware.Auth(h.RenderUser, d))
	e.GET("/about", middleware.Auth(h.About, d))
	e.POST("/logout", h.Logout)
	e.GET("/users", h.AllUsers)
	e.POST("/register", h.RegisterUser)
	e.POST("/follow/:uuid", middleware.Auth(h.Follow, d))
	e.POST("/login", h.Login)
}
