package db

import (
	"htmx/types"

	"github.com/google/uuid"
)

type MockDB struct{}

func CreateMockDB() DB {
	return &MockDB{}
}

var users map[string]user = map[string]user{}

func (mdb MockDB) GetAllUsers() []types.User {
	result := []types.User{}
	for _, u := range users {
		result = append(result, toUser(u))
	}
	return result
}

func (mdb MockDB) InsertUser(username, hashedPassword string, id uuid.UUID) uuid.UUID {
	users[id.String()] = user{username, hashedPassword, id.String()}
	return id
}

func (mdb MockDB) FindUser(id uuid.UUID) (*types.User, string) {
	foundUser, ok := users[id.String()]
	if !ok {
		return nil, ""
	}
	u := toUser(foundUser)
	return &u, foundUser.Password
}

func (mdb MockDB) FindUserByName(name string) (*types.User, string) {
	for _, u := range users {
		if u.Username == name {
			found := toUser(u)
			return &found, u.Password
		}
	}
	return nil, ""
}
