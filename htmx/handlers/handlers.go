package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type user struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

var users map[string]user

// TODO: JWT

func Index(c echo.Context) error {
	fmt.Println(c.Cookies())
	// If auth cookie is valid, return logged in page (with cookie in header?)
	return c.Render(http.StatusOK, "index", "")
}

func RegisterUser(c echo.Context) error {
	u := user{}

	// TODO: Hash password
	users[u.Username] = u

	if err := c.Bind(u); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	data := map[string]string{"Username": u.Username}
	return c.Render(http.StatusOK, "logged-in", data)
}

func Login(c echo.Context) error {
	fmt.Println("Body", c.Request().Body)
	return c.HTML(http.StatusOK, "Logged in")
}
