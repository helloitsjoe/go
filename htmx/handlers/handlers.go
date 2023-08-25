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

type User struct {
	Username string
}

type ctx map[string]interface{}

var users = map[string]user{}
var Users = map[string]User{}

// TODO: JWT
// TODO: SQLite

func SeedUsers() {
	u := [3]string{"Alice", "Bob", "Carl"}

	for _, name := range u {
		Users[name] = User{name}
	}
}

func Index(c echo.Context) error {
	// fmt.Println(c.Cookies())
	// If auth cookie is valid, return logged in page (with cookie in header?)
	data := ctx{"Register": "true", "Users": Users}
	return c.Render(http.StatusOK, "index", data)
}

func RegisterUser(c echo.Context) error {
	fmt.Println(c.Request().FormValue("password"))
	u := user{}

	if err := c.Bind(&u); err != nil {
		fmt.Println(err)
		return c.String(http.StatusBadRequest, "bad request")
	}

	// TODO: Hash password
	users[u.Username] = u

	newUser := User{u.Username}
	Users[u.Username] = newUser

	data := ctx{"NewUser": newUser, "Users": Users}

	return c.Render(http.StatusOK, "logged-in", data)
}

func Login(c echo.Context) error {
	fmt.Println("Body", c.Request().Body)
	return c.HTML(http.StatusOK, "Logged in")
}

func RenderRegister(c echo.Context) error {
	return c.Render(http.StatusOK, "index", ctx{"Register": true})
}

func RenderLogin(c echo.Context) error {
	return c.Render(http.StatusOK, "index", ctx{"Register": false})
}

func AllUsers(c echo.Context) error {
	return c.Render(http.StatusOK, "users", Users)
}
