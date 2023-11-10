package db

import (
	"htmx/types"

	"github.com/hashicorp/go-memdb"
)

type DB interface {
	InsertUser(username, hashedPassword string, id string) string
	FindUser(id string) (*types.User, string)
	FindUserByName(name string) (*types.User, string)
	GetAllUsers() []types.User
	FollowUser(follower, followee string)
	IsFollowing(follower, followee string) bool
	GetFollowers(followers []string) []*types.User
	GetFollowing(following []string) []*types.User
}

// UUIDs need to be converted to strings for memdb,
// otherwise schema gives a length mismatch error
type user struct {
	Username  string
	Password  string
	Followers []string
	Following []string
	UUID      string
}

type MemDB struct {
	db *memdb.MemDB
}

func CreateDB() DB {
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"users": {
				Name: "users",
				Indexes: map[string]*memdb.IndexSchema{
					"username": {
						Name:    "username",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Username"},
					},
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.UUIDFieldIndex{Field: "UUID"},
					},
					"password": {
						Name:    "password",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Password"},
					},
					"followers": {
						Name:    "followers",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Followers"},
					},
					"following": {
						Name:    "following",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Following"},
					},
				},
			},
			// "posts": {
			// 	Name: "posts",
			// 	Indexes: map[string]*memdb.IndexSchema{
			// 		"content": {
			// 			Name:    "content",
			// 			Unique:  false,
			// 			Indexer: &memdb.StringFieldIndex{Field: "Content"},
			// 		},
			// 		"id": {
			// 			Name:    "id",
			// 			Unique:  true,
			// 			Indexer: &memdb.UUIDFieldIndex{Field: "UUID"},
			// 		},
			// 	},
			// },
		},
	}

	db, err := memdb.NewMemDB(schema)
	if err != nil {
		panic(err)
	}

	d := &MemDB{db}

	return d
}
