package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ctx map[string]interface{}

var users = map[string]user{}
var Users = map[string]User{}

// TODO: JWT
// TODO: SQLite

// TODO: Move these to DB
type user struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

type User struct {
	Username string
}

type Handlers struct {
	db string
}

func New(db string) Handlers {

	return Handlers{db}
}

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
	return c.Render(http.StatusOK, "index.html", data)
}

func RegisterUser(c echo.Context) error {
	fmt.Println(c.Request().FormValue("password"))
	u := user{}

	if err := c.Bind(&u); err != nil {
		fmt.Println(err)
		return c.String(http.StatusBadRequest, "bad request")
	}

	if u.Username == "" || u.Password == "" {
		fmt.Println("Name and password must be provided")
		c.Response().Header().Set("HX-Retarget", "#error")
		c.Response().Header().Set("HX-Reswap", "innerHTML")
		return c.Render(http.StatusOK, "error.html", ctx{"Error": "Name and password must be provided"})
	}

	// TODO: Hash password
	users[u.Username] = u

	newUser := User{u.Username}
	Users[u.Username] = newUser

	data := ctx{"User": newUser, "Users": Users}

	return c.Render(http.StatusOK, "logged-in.html", data)
}

func Login(c echo.Context) error {
	fmt.Println("Body", c.Request().Body)
	return c.HTML(http.StatusOK, "Logged in")
}

func RenderRegister(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", ctx{"Register": true, "Users": Users})
}

func RenderLogin(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", ctx{"Register": false, "Users": Users})
}

func AllUsers(c echo.Context) error {
	if c.QueryParam("format") == "json" {
		return c.JSON(http.StatusOK, Users)
	}

	return c.Render(http.StatusOK, "users.html", Users)
}

func About(c echo.Context) error {
	return c.Render(http.StatusOK, "about.html", ctx{"Users": Users})
}
