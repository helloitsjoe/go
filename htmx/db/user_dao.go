package db

import (
	"fmt"
	"htmx/types"

	"github.com/google/uuid"
)

func toUser(u user) types.User {
	id, err := uuid.Parse(u.UUID)
	if err != nil {
		panic(err)
	}
	return types.User{Username: u.Username, UUID: id, Followers: u.Followers, Following: u.Following}
}

func (d MemDB) InsertUser(username, hashedPassword string, id uuid.UUID) uuid.UUID {
	u := user{username, hashedPassword, []uuid.UUID{}, []uuid.UUID{}, id.String()}
	txn := d.db.Txn(true)
	if err := txn.Insert("users", u); err != nil {
		panic(err)
	}

	txn.Commit()

	return id
}

func (db MemDB) FindUser(id uuid.UUID) (*types.User, string) {
	txn := db.db.Txn(false)
	defer txn.Abort()
	all := db.GetAllUsers()
	fmt.Println(all)
	u, err := txn.First("users", "id", id.String())
	if err != nil {
		panic(err)
	}

	// TODO: Maybe return error if not found?
	if u == nil {
		return nil, ""
	}

	foundUser := u.(user)
	result := toUser(foundUser)

	return &result, foundUser.Password
}

func (db MemDB) FindUserByName(name string) (*types.User, string) {
	txn := db.db.Txn(false)
	defer txn.Abort()
	all := db.GetAllUsers()
	fmt.Println(all)
	u, err := txn.First("users", "username", name)
	if err != nil {
		panic(err)
	}

	// TODO: Maybe return error if not found?
	if u == nil {
		return nil, ""
	}

	foundUser := u.(user)
	result := toUser(foundUser)

	return &result, foundUser.Password
}

func (db MemDB) GetAllUsers() []types.User {
	txn := db.db.Txn(false)
	defer txn.Abort()
	it, err := txn.Get("users", "username")
	// TODO: What if not found?
	if err != nil {
		panic(err)
	}

	u := []types.User{}
	for obj := it.Next(); obj != nil; obj = it.Next() {
		foundUser := toUser(obj.(user))
		u = append(u, foundUser)
	}

	return u
}

func (db MemDB) FollowUser(follower, followee uuid.UUID) {
	txn := db.db.Txn(true)
	defer txn.Abort()
	u, _ := db.FindUser(follower)
	u.Following = append(u.Following, followee)
	if err := txn.Insert("users", u); err != nil {
		panic(err)
	}
	txn.Commit()
}
