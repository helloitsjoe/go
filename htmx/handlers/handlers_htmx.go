package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
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
	Id       int
}

type Handlers struct {
	db string
}

var uniqueId = 0

func New(db string) Handlers {

	return Handlers{db}
}

func SeedUsers() {
	u := [3]string{"Alice", "Bob", "Carl"}

	for _, name := range u {
		Users[name] = User{name, uniqueId}
		users[name] = user{name, "password"}
	}
}

func Index(c echo.Context) error {
	if c.Get("user") != nil {
		data := ctx{"Users": Users, "User": Users[c.Get("user").(string)]}
		return c.Render(http.StatusOK, "index.html", data)
	}
	// fmt.Println(c.Cookies())
	// If auth cookie is valid, return logged in page
	data := ctx{"Register": "true", "Users": Users}
	return c.Render(http.StatusOK, "index.html", data)
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
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

	hashed, err := hashPassword(u.Password)
	if err != nil {
		fmt.Println("Error hashing password", err)
		return c.String(http.StatusInternalServerError, "internal server error")
	}

	u.Password = hashed
	users[u.Username] = u

	newUser := User{u.Username, uniqueId}
	Users[u.Username] = newUser

	data := ctx{"User": newUser, "Users": Users}

	return c.Render(http.StatusOK, "logged-in.html", data)
}

func Login(c echo.Context) error {
	fmt.Println("Body", c.Request().Body)
	fmt.Println("Password", c.Request().FormValue("password"))

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

	hashed, err := hashPassword(u.Password)
	if err != nil {
		fmt.Println("Error hashing password", err)
		return c.String(http.StatusInternalServerError, "internal server error")
	}

	user := users[u.Username]
	if !checkPasswordHash(user.Password, hashed) {
		fmt.Println("Incorrect password")
		c.Response().Header().Set("HX-Retarget", "#error")
		c.Response().Header().Set("HX-Reswap", "innerHTML")
		return c.Render(http.StatusOK, "error.html", ctx{"Error": "Incorrect password"})
	}

	loggedInUser := Users[u.Username]

	data := ctx{"User": loggedInUser, "Users": Users}
	fmt.Println(data)

	// TODO: Return better response, clear error message
	return c.Render(http.StatusOK, "logged-in.html", data)
}

func LoggedIn(c echo.Context) error {
	return c.Render(http.StatusOK, "logged-in.html", ctx{"Users": Users})
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
