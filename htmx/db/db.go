package db

import (
	"fmt"
	"htmx/types"

	"github.com/google/uuid"
	"github.com/hashicorp/go-memdb"
)

type DB interface {
	InsertUser(username, hashedPassword string, id uuid.UUID) uuid.UUID
	FindUser(id uuid.UUID) (*types.User, string)
	FindUserByName(name string) (*types.User, string)
	GetAllUsers() []types.User
}

type user struct {
	Username string
	Password string
	UUID     string
}

type MemDB struct {
	db *memdb.MemDB
}

func toUser(u user) types.User {
	id, err := uuid.Parse(u.UUID)
	if err != nil {
		panic(err)
	}
	return types.User{Username: u.Username, UUID: id}
}

func CreateDB() DB {
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

	d := &MemDB{db}

	return d
}

func (d MemDB) InsertUser(username, hashedPassword string, id uuid.UUID) uuid.UUID {
	u := user{username, hashedPassword, id.String()}
	txn := d.db.Txn(true)
	if err := txn.Insert("user", u); err != nil {
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
	u, err := txn.First("user", "id", id.String())
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
	u, err := txn.First("user", "username", name)
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
	it, err := txn.Get("user", "username")
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
