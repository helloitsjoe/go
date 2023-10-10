package db

import (
	"os/user"

	"github.com/google/uuid"
	"github.com/hashicorp/go-memdb"
)

// type user struct {
// 	Username string `form:"username"`
// 	Password string `form:"password"`
// }

// type User struct {
// 	Username string
// }

type DB struct {
	db *memdb.MemDB
}

func CreateDB() *memdb.MemDB {
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"user": {
				Name: "user",
				Indexes: map[string]*memdb.IndexSchema{
					"username": {
						Name:    "username",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Username"},
					},
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Id"},
					},
					"password": {
						Name:    "password",
						Indexer: &memdb.StringFieldIndex{Field: "Password"},
					},
				},
			},
		},
	}

	db, err := memdb.NewMemDB(schema)
	if err != nil {
		panic(err)
	}

	return db
}

func (db DB) AddUser(username, password string) user.User {
	id := uuid.New()

	u := user.User{username, id}
	return u
}
