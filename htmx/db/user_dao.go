package db

import (
	"errors"
	"htmx/types"

	"github.com/hashicorp/go-memdb"
)

func toUser(u user) types.User {
	return types.User{Username: u.Username, UUID: u.UUID, Followers: u.Followers, Following: u.Following}
}

func (d MemDB) InsertUser(username, hashedPassword string, id string) string {
	u := user{username, hashedPassword, []string{}, []string{}, id}
	txn := d.db.Txn(true)
	if err := txn.Insert("users", u); err != nil {
		panic(err)
	}

	txn.Commit()

	return id
}

func (db MemDB) FindUser(id string) (*types.User, string) {
	txn := db.db.Txn(false)
	defer txn.Abort()
	u, err := txn.First("users", "id", id)
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

func findUser(txn *memdb.Txn, id string) user {
	u, err := txn.First("users", "id", id)
	if err != nil {
		panic(err)
	}

	if u == nil {
		panic(errors.New("User not found"))
	}

	// TODO: MemDB Generics?
	return u.(user)
}

func (db MemDB) FollowUser(followerId, followeeId string) {
	txn := db.db.Txn(true)
	defer txn.Abort()
	follower := findUser(txn, followerId)
	followee := findUser(txn, followeeId)

	// TODO: Separate table for followers/following?
	follower.Following = append(follower.Following, followeeId)
	followee.Followers = append(followee.Followers, followerId)
	if err := txn.Insert("users", follower); err != nil {
		panic(err)
	}
	if err := txn.Insert("users", followee); err != nil {
		panic(err)
	}
	txn.Commit()
}

func (db MemDB) IsFollowing(followerId, followeeId string) bool {
	txn := db.db.Txn(false)
	defer txn.Abort()
	follower := findUser(txn, followerId)
	followee := findUser(txn, followeeId)

	for _, f := range follower.Following {
		if f == followeeId {
			return true
		}
	}
	for _, f := range followee.Followers {
		if f == followerId {
			return true
		}
	}
	return false
}

func (db MemDB) GetFollowers(followerIds []string) []*types.User {
	f := []*types.User{}
	for _, followerId := range followerIds {
		follower, _ := db.FindUser(followerId)
		f = append(f, follower)
	}
	return f
}

func (db MemDB) GetFollowing(followingIds []string) []*types.User {
	f := []*types.User{}
	for _, followingId := range followingIds {
		following, _ := db.FindUser(followingId)
		f = append(f, following)
	}
	return f
}

// TODO: Unfollow user
