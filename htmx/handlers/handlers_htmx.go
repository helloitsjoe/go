package handlers

import (
	"fmt"
	"htmx/user"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ctx map[string]interface{}

type Handlers struct {
	db string
}

func New(db string) Handlers {
	return Handlers{db}
}

func Index(c echo.Context) error {
	if c.Get("user") != nil {
		data := ctx{"Users": user.Users, "User": user.Users[c.Get("user").(string)]}
		return c.Render(http.StatusOK, "index.html", data)
	}
	// fmt.Println(c.Cookies())
	// If auth cookie is valid, return logged in page
	data := ctx{"Register": "true", "Users": user.Users}
	return c.Render(http.StatusOK, "index.html", data)
}

func RegisterUser(c echo.Context) error {
	fmt.Println(c.Request().FormValue("password"))
	username := c.Request().FormValue("username")
	password := c.Request().FormValue("password")

	newUser, err := user.AddUser(c, username, password)
	if err != nil {
		if err.Error() == "Bad request" {
			return c.String(http.StatusBadRequest, "bad request")
		} else if err.Error() == "Error hashing password" {
			return c.String(http.StatusInternalServerError, "internal server error")
		}
		fmt.Println("Error adding user", err)
		c.Response().Header().Set("HX-Retarget", "#error")
		c.Response().Header().Set("HX-Reswap", "innerHTML")
		return c.Render(http.StatusOK, "error.html", ctx{"Error": err.Error()})
	}

	data := ctx{"User": newUser, "Users": user.Users}

	return c.Render(http.StatusOK, "logged-in.html", data)
}

func Login(c echo.Context) error {
	fmt.Println("Body", c.Request().Body)
	fmt.Println("Password", c.Request().FormValue("password"))
	username := c.Request().FormValue("username")
	password := c.Request().FormValue("password")

	loggedInUser, err := user.Login(c, username, password)
	if err != nil {
		if err.Error() == "Bad request" {
			return c.String(http.StatusBadRequest, "bad request")
		} else if err.Error() == "Error hashing password" {
			return c.String(http.StatusInternalServerError, "internal server error")
		}
		fmt.Println("Error adding user", err)
		c.Response().Header().Set("HX-Retarget", "#error")
		c.Response().Header().Set("HX-Reswap", "innerHTML")
		return c.Render(http.StatusOK, "error.html", ctx{"Error": err.Error()})
	}

	data := ctx{"User": loggedInUser, "Users": user.Users}
	fmt.Println(data)

	return c.Render(http.StatusOK, "logged-in.html", data)
}

func LoggedIn(c echo.Context) error {
	return c.Render(http.StatusOK, "logged-in.html", ctx{"Users": user.Users})
}

func RenderRegister(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", ctx{"Register": true, "Users": user.Users})
}

func RenderLogin(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", ctx{"Register": false, "Users": user.Users})
}

func AllUsers(c echo.Context) error {
	if c.QueryParam("format") == "json" {
		return c.JSON(http.StatusOK, user.Users)
	}

	return c.Render(http.StatusOK, "users.html", user.Users)
}

func About(c echo.Context) error {
	return c.Render(http.StatusOK, "about.html", ctx{"Users": user.Users})
}
