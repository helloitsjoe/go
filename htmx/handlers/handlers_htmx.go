package handlers

import (
	"fmt"
	"htmx/db"
	"htmx/types"
	"htmx/user"
	"net/http"
	"sort"
	"strconv"
	"strings"
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

func getLoginCookie(user *types.User) *http.Cookie {
	return &http.Cookie{Name: "uuid", Value: user.UUID, HttpOnly: true, MaxAge: 10 * 60}
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

func (h Handlers) redirectHome(c echo.Context) error {
	c.Response().Header().Set("HX-Redirect", "/")
	users := h.db.GetAllUsers()
	return c.Render(http.StatusOK, "index.html", ctx{"Register": false, "Users": users})
}

func (h Handlers) Index(c echo.Context) error {
	loggedInUser, ok := c.Get("user").(*types.User)
	users := h.db.GetAllUsers()

	if ok {
		fmt.Println("followers", loggedInUser.Followers)
		data := ctx{"Users": users, "User": loggedInUser}
		return c.Render(http.StatusOK, "index.html", data)
	}

	fmt.Print("Not logged in")
	data := ctx{"Register": c.Path() != "/login", "Users": users}
	return c.Render(http.StatusOK, "index.html", data)
}

func (h Handlers) RegisterUser(c echo.Context) error {
	sleep := getSleep(c.Request().FormValue("sleep"))
	time.Sleep(sleep * time.Second)
	username := c.Request().FormValue("username")
	password := c.Request().FormValue("password")

	users := user.GetUsers(h.db)
	newUser, err := user.AddUser(h.db, strings.ToLower(username), password)

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

	loggedInUser, err := user.Login(h.db, strings.ToLower(username), password)
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
	loggedInUser, ok := c.Get("user").(*types.User)
	users := h.db.GetAllUsers()

	if ok {
		return c.Render(http.StatusOK, "about.html", ctx{"Users": users, "User": loggedInUser})
	}

	return c.Render(http.StatusOK, "about.html", ctx{"Users": users})
}

func (h Handlers) Logout(c echo.Context) error {
	c.SetCookie(&http.Cookie{Name: "uuid", Value: "", HttpOnly: true, MaxAge: -1})
	return h.redirectHome(c)
}

func (h Handlers) RenderFollowers(c echo.Context) error {
	loggedInUser, ok := c.Get("user").(*types.User)

	if ok {
		followers := user.GetFollowers(h.db, loggedInUser.Followers)
		data := ctx{"User": loggedInUser, "Followers": followers}

		fmt.Println("data", data)
		return c.Render(http.StatusOK, "followers.html", data)
	}
	return h.redirectHome(c)
}

func (h Handlers) RenderFollowing(c echo.Context) error {
	loggedInUser, ok := c.Get("user").(*types.User)

	if ok {
		following := user.GetFollowing(h.db, loggedInUser.Following)
		data := ctx{"User": loggedInUser, "Following": following}
		return c.Render(http.StatusOK, "following.html", data)
	}
	return h.redirectHome(c)
}

func (h Handlers) RenderUser(c echo.Context) error {
	loggedInUser := c.Get("user")

	if loggedInUser != nil {
		loggedInUser = loggedInUser.(*types.User)
	}

	user := c.Param("username")
	if user == "" {
		return h.redirectHome(c)
	}
	fmt.Println("user", user)

	u, _ := h.db.FindUserByName(user)

	booksCheckedOut := []string{"The Hobbit"}
	booksAvailable := []string{"Dune"}

	allUsers := h.db.GetAllUsers()

	fmt.Println("loggedInUser", loggedInUser)

	data := ctx{"User": loggedInUser, "TargetUser": u, "BooksCheckedOut": booksCheckedOut, "BooksAvailable": booksAvailable, "Users": allUsers}
	return c.Render(http.StatusOK, "user.html", data)
}

func (h Handlers) Follow(c echo.Context) error {
	loggedInUser, ok := c.Get("user").(*types.User)
	if !ok {
		// TODO: if loggedInUser == nil
		fmt.Println("No user!")
	}

	targetId := c.Param("uuid")

	user.Follow(h.db, loggedInUser.UUID, targetId)

	// TODO: Toggle "Unfollowing"
	return c.HTML(http.StatusOK, "<button>Following</button>")
}

// TODO: Borrow, return
