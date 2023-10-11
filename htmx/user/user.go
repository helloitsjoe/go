package user

import (
	"errors"
	"fmt"
	"htmx/db"
	"htmx/types"

	"github.com/google/uuid"
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

// var users = map[string]*user{}
var Users = map[string]*types.User{}

// TODO: JWT

func NewUser(username string) *types.User {
	u := &types.User{}
	u.Username = username
	u.UUID = uuid.New()
	return u
}

func SeedUsers(db *db.DB) {
	u := [3]string{"Alice", "Bob", "Carl"}

	for _, name := range u {
		n := NewUser(name)
		fmt.Println(n)
		db.InsertUser(n.Username, "bar", n.UUID)
	}
}

func AddUser(c echo.Context, db *db.DB, name, password string) (*types.User, error) {
	u := NewUser(name)

	// TODO: Extract Context out of this function
	if err := c.Bind(u); err != nil {
		fmt.Println(err)
		return nil, errors.New("Bad request")
	}

	if u.Username == "" || password == "" {
		fmt.Println("Name and password must be provided")
		return nil, errors.New("Name and password must be provided")
	}

	hashed, err := hashPassword(password)
	if err != nil {
		fmt.Println("Error hashing password", err)
		return nil, errors.New("Error hashing password")
	}

	db.InsertUser(u.Username, hashed, u.UUID)

	return u, nil
}

func Login(c echo.Context, db *db.DB, name, password string) (*types.User, error) {
	u := NewUser(name)

	// TODO: Extract Context out of this function
	if err := c.Bind(u); err != nil {
		fmt.Println(err)
		return nil, errors.New("Bad request")
	}

	if u.Username == "" || password == "" {
		fmt.Println("Name and password must be provided")
		return nil, errors.New("Name and password must be provided")
	}

	hashed, err := hashPassword(password)
	if err != nil {
		fmt.Println("Error hashing password", err)
		return nil, errors.New("Error hashing password")
	}

	// user := users[u.Username]
	newUser, userHashed := db.FindUser(u.UUID)

	if !checkPasswordHash(userHashed, hashed) {
		fmt.Println("Incorrect password")
		return nil, errors.New("Incorrect password")
	}

	// loggedInUser := Users[u.Username]

	return newUser, nil
}

func GetUsers(db *db.DB) []types.User {
	users := db.GetAllUsers()

	return users
}
