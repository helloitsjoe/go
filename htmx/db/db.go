package db

import (
	"fmt"
	"htmx/types"
	"reflect"

	"github.com/google/uuid"
	"github.com/hashicorp/go-memdb"
)

type user struct {
	Username string
	Password string
	UUID     string
}

// type User struct {
// 	Username string
// }

type DB struct {
	db *memdb.MemDB
}

func CreateDB() *DB {
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
						Indexer: &memdb.UUIDFieldIndex{Field: "UUID"},
					},
					"password": {
						Name:    "password",
						Unique:  false,
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

	d := &DB{db}

	return d
}

func (d DB) InsertUser(username, hashedPassword string, id uuid.UUID) uuid.UUID {
	u := user{username, hashedPassword, id.String()}
	txn := d.db.Txn(true)
	if err := txn.Insert("user", u); err != nil {
		panic(err)
	}

	txn.Commit()

	return id
}

func (db DB) FindUser(id uuid.UUID) (*types.User, string) {
	txn := db.db.Txn(false)
	defer txn.Abort()
	u, err := txn.First("user", "id", id.String())
	// TODO: What if not found?
	if err != nil {
		panic(err)
	}

	foundUser := u.(user)
	result := &types.User{Username: foundUser.Username, UUID: id}

	return result, foundUser.Password
}

func (db DB) GetAllUsers() []types.User {
	txn := db.db.Txn(false)
	defer txn.Abort()
	it, err := txn.Get("user", "id")
	// TODO: What if not found?
	if err != nil {
		panic(err)
	}

	u := []types.User{}
	for obj := it.Next(); obj != nil; obj = it.Next() {
		fmt.Println("here", reflect.TypeOf(obj.(user).UUID))
		id, err := uuid.Parse(obj.(user).UUID)
		if err != nil {
			panic(err)
		}
		foundUser := types.User{Username: obj.(user).Username, UUID: id}
		// fmt.Println(foundUser)
		u = append(u, foundUser)
		// fmt.Println(u)
	}

	return u
}
