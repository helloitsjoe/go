package handlers

import (
	"fmt"
	"htmx/user"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type ctx map[string]interface{}

type Handlers struct {
	db string
}

func New(db string) Handlers {
	return Handlers{db}
}

func checkLoggedIn(username string, usernameExists bool) (*user.User, bool) {
	if usernameExists {
		validUser, userExists := user.Users[username]

		if userExists {
			return &validUser, true
		}
	}

	return nil, false
}

func getSleep(reqSleep string) time.Duration {
	sleep := 0
	if reqSleep, err := strconv.Atoi(reqSleep); err == nil {
		sleep = reqSleep
	} else {
		fmt.Println("Error parsing sleep-seconds:", err)
	}
	return time.Duration(sleep)
}

func Index(c echo.Context) error {
	username, ok := c.Get("username").(string)
	loggedInUser, exists := checkLoggedIn(username, ok)

	if exists {
		data := ctx{"Users": user.Users, "User": loggedInUser}
		return c.Render(http.StatusOK, "index.html", data)
	}

	data := ctx{"Register": "true", "Users": user.Users}
	return c.Render(http.StatusOK, "index.html", data)
}

func RegisterUser(c echo.Context) error {
	sleep := getSleep(c.Request().FormValue("sleep"))
	time.Sleep(sleep * time.Second)
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

	c.SetCookie(&http.Cookie{Name: "username", Value: newUser.Username, HttpOnly: true, MaxAge: 10 * 60})

	data := ctx{"User": newUser, "Users": user.Users}

	return c.Render(http.StatusOK, "logged-in.html", data)
}

func Login(c echo.Context) error {
	sleep := getSleep(c.Request().FormValue("sleep"))
	time.Sleep(time.Duration(sleep) * time.Second)
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

	c.SetCookie(&http.Cookie{Name: "username", Value: loggedInUser.Username, HttpOnly: true, MaxAge: 10 * 60})

	return c.Render(http.StatusOK, "logged-in.html", data)
}

func LoggedIn(c echo.Context) error {
	return c.Render(http.StatusOK, "logged-in.html", ctx{"Users": user.Users})
}

func RenderRegister(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", ctx{"Register": true, "Users": user.Users})
}

func RenderLogin(c echo.Context) error {
	authUsername, ok := c.Get("username").(string)
	_, exists := checkLoggedIn(authUsername, ok)

	if exists {
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}

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

func Logout(c echo.Context) error {
	c.SetCookie(&http.Cookie{Name: "username", Value: "", HttpOnly: true, MaxAge: 0})
	c.Response().Header().Set("HX-Redirect", "/")
	return c.Render(http.StatusOK, "index.html", ctx{"Register": false, "Users": user.Users})
}
