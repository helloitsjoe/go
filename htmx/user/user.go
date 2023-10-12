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

func hashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		panic(err)
	}
	return string(bytes)
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

	users := map[string]*types.User{}

	for _, name := range u {
		n := NewUser(name)
		users[name] = n
		fmt.Println("inserting user", n)
		before := db.GetAllUsers()
		fmt.Println("all users before", before)
		p := hashPassword("bar")
		id := db.InsertUser(n.Username, p, n.UUID)
		fmt.Println("id", id)
		after := db.GetAllUsers()
		fmt.Println("all users after", after)
	}
	uu := users["Alice"].UUID
	fmt.Println("uu", uu)
	a, _ := db.FindUser(uu)
	fmt.Println("Alice", a)
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

	hashed := hashPassword(password)

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

	hashed := hashPassword(password)

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
