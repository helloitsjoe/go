package db

import (
	"errors"
	"htmx/types"

	"github.com/google/uuid"
	"github.com/hashicorp/go-memdb"
)

func toUser(u user) types.User {
	id, err := uuid.Parse(u.UUID)
	if err != nil {
		panic(err)
	}
	return types.User{Username: u.Username, UUID: id, Followers: u.Followers, Following: u.Following}
}

func (d MemDB) InsertUser(username, hashedPassword string, id uuid.UUID) uuid.UUID {
	u := user{username, hashedPassword, []string{}, []string{}, id.String()}
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

func (db MemDB) FollowUser(followerId, followeeId uuid.UUID) {
	txn := db.db.Txn(true)
	defer txn.Abort()
	follower := findUser(txn, followerId.String())
	followee := findUser(txn, followeeId.String())

	// TODO: Separate table for followers/following?
	follower.Following = append(follower.Following, followeeId.String())
	followee.Followers = append(followee.Followers, followerId.String())
	if err := txn.Insert("users", follower); err != nil {
		panic(err)
	}
	if err := txn.Insert("users", followee); err != nil {
		panic(err)
	}
	txn.Commit()
}

func (db MemDB) IsFollowing(followerId, followeeId uuid.UUID) bool {
	txn := db.db.Txn(false)
	defer txn.Abort()
	follower := findUser(txn, followerId.String())
	followee := findUser(txn, followeeId.String())

	for _, f := range follower.Following {
		if f == followeeId.String() {
			return true
		}
	}
	for _, f := range followee.Followers {
		if f == followerId.String() {
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
