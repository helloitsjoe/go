package user

import (
	"errors"
	"fmt"
	"htmx/db"
	"htmx/types"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func checkPasswordHash(plaintext, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plaintext))
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

func AddUser(db *db.DB, name, password string) (*types.User, error) {
	u := NewUser(name)

	if u.Username == "" || password == "" {
		fmt.Println("Name and password must be provided")
		return nil, errors.New("Name and password must be provided")
	}

	hashed := hashPassword(password)

	db.InsertUser(u.Username, hashed, u.UUID)

	return u, nil
}

func Login(db *db.DB, name, password string) (*types.User, error) {
	// TODO: Separate function for finding hashed password?
	u, userHashed := db.FindUserByName(name)

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
