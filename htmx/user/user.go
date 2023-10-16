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
		p := hashPassword("bar")
		fmt.Println("p", p)
		db.InsertUser(n.Username, p, n.UUID)
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

	hashed := hashPassword(password)

	db.InsertUser(u.Username, hashed, u.UUID)

	return u, nil
}

func Login(c echo.Context, db *db.DB, name, password string) (*types.User, error) {
	// TODO: Separate function for finding hashed password?
	u, userHashed := db.FindUserByName(name)

	// TODO: Extract Context out of this function
	// if err := c.Bind(u); err != nil {
	// 	fmt.Println(err)
	// 	return nil, errors.New("Bad request")
	// }

	fmt.Println("in Login", password)

	if u.Username == "" || password == "" {
		fmt.Println("Name and password must be provided")
		return nil, errors.New("Name and password must be provided")
	}

	if !checkPasswordHash(password, userHashed) {
		fmt.Println("Incorrect password")
		return nil, errors.New("Incorrect password")
	}

	return u, nil
}

func GetUsers(db *db.DB) []types.User {
	users := db.GetAllUsers()

	return users
}
