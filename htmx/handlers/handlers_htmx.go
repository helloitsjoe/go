package handlers

import (
	"fmt"
	"htmx/db"
	"htmx/types"
	"htmx/user"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type ctx map[string]interface{}

type Handlers struct {
	db db.DB
}

func NewHandlers(db db.DB) Handlers {
	return Handlers{db}
}

func checkLoggedIn(id string, idExists bool, db db.DB) (*types.User, bool) {
	if !idExists {
		return nil, false
	}

	validUser, _ := db.FindUser(id)

	if validUser != nil {
		return validUser, true
	}

	return nil, false
}

func getLoginCookie(user *types.User) *http.Cookie {
	return &http.Cookie{Name: "uuid", Value: user.UUID.String(), HttpOnly: true, MaxAge: 10 * 60}
}

func getSleep(reqSleep string) time.Duration {
	sleep := 0
	if reqSleep, err := strconv.Atoi(reqSleep); err == nil {
		sleep = reqSleep
	} else {
		fmt.Println("Error parsing sleep-seconds:", err)
		fmt.Println("Continuing, sleeping for", sleep, "seconds")
	}
	return time.Duration(sleep)
}

func (h Handlers) Index(c echo.Context) error {
	id, ok := c.Get("uuid").(string)
	loggedInUser, isLoggedIn := checkLoggedIn(id, ok, h.db)
	users := h.db.GetAllUsers()

	if isLoggedIn {
		fmt.Println(loggedInUser.Followers)
		data := ctx{"Users": users, "User": loggedInUser}
		return c.Render(http.StatusOK, "index.html", data)
	}

	data := ctx{"Register": c.Path() != "/login", "Users": users}
	return c.Render(http.StatusOK, "index.html", data)
}

func (h Handlers) RegisterUser(c echo.Context) error {
	sleep := getSleep(c.Request().FormValue("sleep"))
	time.Sleep(sleep * time.Second)
	username := c.Request().FormValue("username")
	password := c.Request().FormValue("password")

	users := user.GetUsers(h.db)
	newUser, err := user.AddUser(h.db, username, password)

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

	users = append(users, *newUser)
	c.SetCookie(getLoginCookie(newUser))

	data := ctx{"User": newUser, "Users": users}

	return c.Render(http.StatusOK, "logged-in.html", data)
}

func (h Handlers) Login(c echo.Context) error {
	sleep := getSleep(c.Request().FormValue("sleep"))
	time.Sleep(sleep * time.Second)
	fmt.Println("Body", c.Request().Body)
	fmt.Println("Password", c.Request().FormValue("password"))
	username := c.Request().FormValue("username")
	password := c.Request().FormValue("password")

	loggedInUser, err := user.Login(h.db, username, password)
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

	users := h.db.GetAllUsers()

	c.SetCookie(getLoginCookie(loggedInUser))

	data := ctx{"User": loggedInUser, "Users": users}
	return c.Render(http.StatusOK, "logged-in.html", data)
}

func (h Handlers) LoggedIn(c echo.Context) error {
	users := h.db.GetAllUsers()
	return c.Render(http.StatusOK, "logged-in.html", ctx{"Users": users})
}

func (h Handlers) AllUsers(c echo.Context) error {
	users := h.db.GetAllUsers()
	sort.Slice(users, func(i, j int) bool {
		return users[i].Username < users[j].Username
	})
	if c.QueryParam("format") == "json" {
		return c.JSON(http.StatusOK, users)
	}

	return c.Render(http.StatusOK, "users.html", users)
}

func (h Handlers) About(c echo.Context) error {
	users := h.db.GetAllUsers()
	return c.Render(http.StatusOK, "about.html", ctx{"Users": users})
}

func (h Handlers) Logout(c echo.Context) error {
	c.SetCookie(&http.Cookie{Name: "uuid", Value: "", HttpOnly: true, MaxAge: -1})
	c.Response().Header().Set("HX-Redirect", "/")
	users := h.db.GetAllUsers()
	return c.Render(http.StatusOK, "index.html", ctx{"Register": false, "Users": users})
}

func (h Handlers) RenderFollowers(c echo.Context) error {
	id, ok := c.Get("uuid").(string)
	loggedInUser, isLoggedIn := checkLoggedIn(id, ok, h.db)

	if isLoggedIn {
		followers := user.GetFollowers(h.db, loggedInUser.Followers)
		data := ctx{"User": loggedInUser, "Followers": followers}
		return c.Render(http.StatusOK, "followers.html", data)
	}
	// TODO: DRY
	c.Response().Header().Set("HX-Redirect", "/")
	users := h.db.GetAllUsers()
	return c.Render(http.StatusOK, "index.html", ctx{"Register": false, "Users": users})
}

func (h Handlers) RenderFollowing(c echo.Context) error {
	// TODO: Pass user from auth
	id, ok := c.Get("uuid").(string)
	loggedInUser, isLoggedIn := checkLoggedIn(id, ok, h.db)

	if isLoggedIn {
		following := user.GetFollowing(h.db, loggedInUser.Following)
		data := ctx{"User": loggedInUser, "Following": following}
		return c.Render(http.StatusOK, "following.html", data)
	}
	c.Response().Header().Set("HX-Redirect", "/")
	users := h.db.GetAllUsers()
	return c.Render(http.StatusOK, "index.html", ctx{"Register": false, "Users": users})
}
