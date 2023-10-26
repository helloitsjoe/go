package db

import (
	"htmx/types"

	"github.com/google/uuid"
)

// MockDB is unused, but it's an example of how an interface can be used with multiple implementations in Go

type MockDB struct {
	users map[string]user
}

func CreateMockDB() DB {
	users := map[string]user{}
	return &MockDB{users}
}

func (d MockDB) GetAllUsers() []types.User {
	result := []types.User{}
	for _, u := range d.users {
		result = append(result, toUser(u))
	}
	return result
}

func (d MockDB) InsertUser(username, hashedPassword string, id uuid.UUID) uuid.UUID {
	d.users[id.String()] = user{username, hashedPassword, id.String()}
	return id
}

func (d MockDB) FindUser(id uuid.UUID) (*types.User, string) {
	foundUser, ok := d.users[id.String()]
	if !ok {
		return nil, ""
	}
	u := toUser(foundUser)
	return &u, foundUser.Password
}

func (d MockDB) FindUserByName(name string) (*types.User, string) {
	for _, u := range d.users {
		if u.Username == name {
			found := toUser(u)
			return &found, u.Password
		}
	}
	return nil, ""
}