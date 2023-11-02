package user

import (
	"htmx/db"
	"htmx/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSeedUsers(t *testing.T) {
	d := db.CreateDB()
	SeedUsers(d)
	users := d.GetAllUsers()

	alice, bob, carl := users[0], users[1], users[2]

	assert.Equal(t, alice.Username, "Alice")
	assert.Equal(t, bob.Username, "Bob")
	assert.Equal(t, carl.Username, "Carl")

	assert.NotEqual(t, alice.UUID, bob.UUID)
	assert.NotEqual(t, bob.UUID, carl.UUID)
	assert.NotEqual(t, alice.UUID, carl.UUID)

	assert.Equal(t, alice.Following, []string{bob.UUID.String(), carl.UUID.String()})
}

func TestLoginSuccess(t *testing.T) {
	d := db.CreateDB()
	SeedUsers(d)
	u, err := Login(d, "Alice", "bar")
	assert.IsType(t, u, &types.User{})
	assert.Nil(t, err)
}

func TestLoginNoName(t *testing.T) {
	d := db.CreateDB()
	SeedUsers(d)
	u, err := Login(d, "", "bar")
	assert.Nil(t, u)
	assert.Equal(t, err.Error(), "Name and password must be provided")
}

func TestLoginNoPassword(t *testing.T) {
	d := db.CreateDB()
	SeedUsers(d)
	u, err := Login(d, "Alice", "")
	assert.Nil(t, u)
	assert.Equal(t, err.Error(), "Name and password must be provided")
}

func TestLoginMissingUser(t *testing.T) {
	d := db.CreateDB()
	SeedUsers(d)
	u, err := Login(d, "Nobody", "foo")
	assert.Nil(t, u)
	assert.Equal(t, err.Error(), "Incorrect login credentials")
}

func TestLoginIncorrectPassword(t *testing.T) {
	d := db.CreateDB()
	SeedUsers(d)
	u, err := Login(d, "Alice", "dunno")
	assert.Nil(t, u)
	assert.Equal(t, err.Error(), "Incorrect login credentials")
}
