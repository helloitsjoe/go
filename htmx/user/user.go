package user

import (
	"errors"
	"fmt"
	"htmx/db"
	"htmx/types"
	"strings"

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
	u.Followers = []string{}
	u.Following = []string{}
	return u
}

func SeedUsers(d db.DB) {
	u := [3]string{"Alice", "Bob", "Carl"}
	newUsers := []*types.User{}

	for _, name := range u {
		n := NewUser(name)
		p := hashPassword("bar")
		newUsers = append(newUsers, n)
		d.InsertUser(n.Username, p, n.UUID)
	}

	// TODO: Convert all UUIDs to strings at creation time
	alice, bob, carl := newUsers[0], newUsers[1], newUsers[2]

	d.FollowUser(alice.UUID, bob.UUID)
	d.FollowUser(bob.UUID, alice.UUID)
	d.FollowUser(carl.UUID, alice.UUID)
}

func AddUser(db db.DB, name, password string) (*types.User, error) {
	u := NewUser(name)

	// MemDB doesn't enforce uniqueness, so we have to check manually before
	// insesrting: https://github.com/hashicorp/go-memdb/issues/7
	users := db.GetAllUsers()
	for _, u := range users {
		if strings.EqualFold(u.Username, name) {
			return nil, errors.New("That username is already taken")
		}
	}

	if u.Username == "" || password == "" {
		fmt.Println("Name and password must be provided")
		return nil, errors.New("Name and password must be provided")
	}

	hashed := hashPassword(password)

	db.InsertUser(u.Username, hashed, u.UUID)

	return u, nil
}

func Login(db db.DB, name, password string) (*types.User, error) {
	// TODO: Separate function for finding hashed password?
	u, userHashed := db.FindUserByName(name)
	if name == "" || password == "" {
		fmt.Println("Name and password must be provided")
		return nil, errors.New("Name and password must be provided")
	}

	if u == nil {
		fmt.Println("User not found:", name)
		return nil, errors.New("Incorrect login credentials")
	}

	if !checkPasswordHash(password, userHashed) {
		fmt.Println("Incorrect password")
		return nil, errors.New("Incorrect login credentials")
	}

	return u, nil
}

func GetUsers(db db.DB) []types.User {
	users := db.GetAllUsers()

	return users
}

func Follow(db db.DB, a, b uuid.UUID) {
	if a == b {
		return
	}
	if db.IsFollowing(a, b) {
		return
	}
	db.FollowUser(a, b)
}

func GetFollowers(db db.DB, followers []string) []*types.User {
	return db.GetFollowers(followers)
}
