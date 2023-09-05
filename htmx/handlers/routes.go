package handlers

import (
	"github.com/labstack/echo/v4"
)

func Register(e *echo.Echo) {
	SeedUsers()

	e.GET("/", Index)
	e.GET("/about", About)
	e.GET("/register", RenderRegister)
	e.GET("/login", RenderLogin)
	e.GET("/users", AllUsers)
	e.POST("/register", RegisterUser)
	e.POST("/login", Login)

}
