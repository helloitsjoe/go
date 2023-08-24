package handlers

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/labstack/echo/v4"
)

type user struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var users map[string]user

func Index(c echo.Context) error {
	t, err := template.ParseFiles("static/index.html")
	if err != nil {
		return err
	}
	return t.Execute(c.Response(), nil)
}

func RegisterUser(c echo.Context) error {
	fmt.Println("Body", c.Request().Body)
	return c.HTML(http.StatusOK, "Registered")
}

func Login(c echo.Context) error {
	fmt.Println("Body", c.Request().Body)
	return c.HTML(http.StatusOK, "Logged in")
}
