package db

import (
	"htmx/types"
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

func (d MockDB) InsertUser(username, hashedPassword string, id string) string {
	d.users[id] = user{username, hashedPassword, []string{}, []string{}, id}
	return id
}

func (d MockDB) FindUser(id string) (*types.User, string) {
	foundUser, ok := d.users[id]
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

func (d MockDB) FollowUser(follower, followee string) {
	panic("not implemented")
}

func (d MockDB) UnfollowUser(follower, followee string) {
	panic("not implemented")
}

func (d MockDB) IsFollowing(follower, followee string) bool {
	panic("not implemented")
}

func (d MockDB) GetFollowers(followers []string) []*types.User {
	panic("not implemented")
}

func (d MockDB) GetFollowing(following []string) []*types.User {
	panic("not implemented")
}
