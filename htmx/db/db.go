package db

import (
	"fmt"
	"htmx/types"

	"github.com/google/uuid"
	"github.com/hashicorp/go-memdb"
)

type user struct {
	Username string
	Password string
	UUID     string
}

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
	all := db.GetAllUsers()
	fmt.Println(all)
	u, err := txn.First("user", "id", id.String())
	if err != nil {
		panic(err)
	}

	// TODO: Maybe return error if not found?
	if u == nil {
		return nil, ""
	}

	foundUser := u.(user)
	result := &types.User{Username: foundUser.Username, UUID: id}

	return result, foundUser.Password
}

func (db DB) FindUserByName(name string) (*types.User, string) {
	txn := db.db.Txn(false)
	defer txn.Abort()
	all := db.GetAllUsers()
	fmt.Println(all)
	u, err := txn.First("user", "username", name)
	if err != nil {
		panic(err)
	}

	// TODO: Maybe return error if not found?
	if u == nil {
		return nil, ""
	}

	foundUser := u.(user)
	id, err := uuid.Parse(foundUser.UUID)
	if err != nil {
		panic(err)
	}
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
		id, err := uuid.Parse(obj.(user).UUID)
		if err != nil {
			panic(err)
		}
		foundUser := types.User{Username: obj.(user).Username, UUID: id}
		u = append(u, foundUser)
	}

	return u
}
