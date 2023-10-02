package user

import (
	"errors"
	"fmt"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

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

var uniqueId = 0

func SeedUsers() {
	u := [3]string{"Alice", "Bob", "Carl"}

	for _, name := range u {
		uniqueId++
		Users[name] = User{name, uniqueId}
		users[name] = user{name, "bar"}
	}
}

func AddUser(c echo.Context, name, password string) (*User, error) {
	u := user{}

	// TODO: Extract Context out of this function
	if err := c.Bind(&u); err != nil {
		fmt.Println(err)
		return nil, errors.New("Bad request")
	}

	if u.Username == "" || u.Password == "" {
		fmt.Println("Name and password must be provided")
		return nil, errors.New("Name and password must be provided")
	}

	hashed, err := hashPassword(u.Password)
	if err != nil {
		fmt.Println("Error hashing password", err)
		return nil, errors.New("Error hashing password")
	}

	u.Password = hashed
	users[u.Username] = u

	newUser := User{u.Username, uniqueId}
	Users[u.Username] = newUser

	return &newUser, nil
}

func Login(c echo.Context, name, password string) (*User, error) {
	u := user{}

	// TODO: Extract Context out of this function
	if err := c.Bind(&u); err != nil {
		fmt.Println(err)
		return nil, errors.New("Bad request")
	}

	if u.Username == "" || u.Password == "" {
		fmt.Println("Name and password must be provided")
		return nil, errors.New("Name and password must be provided")
	}

	hashed, err := hashPassword(u.Password)
	if err != nil {
		fmt.Println("Error hashing password", err)
		return nil, errors.New("Error hashing password")
	}

	user := users[u.Username]

	if !checkPasswordHash(user.Password, hashed) {
		fmt.Println("Incorrect password")
		c.Response().Header().Set("HX-Retarget", "#error")
		c.Response().Header().Set("HX-Reswap", "innerHTML")
		return nil, errors.New("Incorrect password")
	}

	loggedInUser := Users[u.Username]

	return &loggedInUser, nil
}
