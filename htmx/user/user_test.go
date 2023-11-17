package user

import (
	"fmt"
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

	assert.Equal(t, alice.Username, "alice")
	assert.Equal(t, bob.Username, "bob")
	assert.Equal(t, carl.Username, "carl")

	assert.NotEqual(t, alice.UUID, bob.UUID)
	assert.NotEqual(t, bob.UUID, carl.UUID)
	assert.NotEqual(t, alice.UUID, carl.UUID)

}

func TestFollowUser(t *testing.T) {
	d := db.CreateDB()
	a, _ := AddUser(d, "alice", "bar")
	b, _ := AddUser(d, "bob", "bar")
	c, _ := AddUser(d, "carl", "bar")

	Follow(d, a.UUID, b.UUID)
	Follow(d, b.UUID, a.UUID)
	Follow(d, c.UUID, a.UUID)

	alice, _ := d.FindUser(a.UUID)
	bob, _ := d.FindUser(b.UUID)
	carl, _ := d.FindUser(c.UUID)

	assert.Equal(t, alice.Following, []string{bob.UUID})
	assert.Equal(t, alice.Followers, []string{bob.UUID, carl.UUID})
	assert.Equal(t, bob.Following, []string{alice.UUID})
	assert.Equal(t, bob.Followers, []string{alice.UUID})
	assert.Equal(t, carl.Following, []string{alice.UUID})
	assert.Equal(t, carl.Followers, []string{})
}

func TestUnfollowUser(t *testing.T) {
	d := db.CreateDB()
	a, _ := AddUser(d, "alice", "bar")
	b, _ := AddUser(d, "bob", "bar")
	c, _ := AddUser(d, "carl", "bar")

	Follow(d, a.UUID, b.UUID)
	Follow(d, b.UUID, a.UUID)
	Follow(d, c.UUID, a.UUID)

	alice, _ := d.FindUser(a.UUID)
	bob, _ := d.FindUser(b.UUID)
	carl, _ := d.FindUser(c.UUID)

	assert.Equal(t, alice.Following, []string{bob.UUID})
	assert.Equal(t, alice.Followers, []string{bob.UUID, carl.UUID})
	assert.Equal(t, bob.Following, []string{alice.UUID})
	assert.Equal(t, bob.Followers, []string{alice.UUID})
	assert.Equal(t, carl.Following, []string{alice.UUID})
	assert.Equal(t, carl.Followers, []string{})

	fmt.Println("Alice", a.UUID)
	Unfollow(d, b.UUID, a.UUID)
	Unfollow(d, c.UUID, a.UUID)

	alice, _ = d.FindUser(a.UUID)
	bob, _ = d.FindUser(b.UUID)
	carl, _ = d.FindUser(c.UUID)

	assert.Equal(t, []string{}, alice.Followers)
	assert.Equal(t, []string{}, bob.Following)
	assert.Equal(t, []string{}, carl.Following)
}

func TestFollowSelfFail(t *testing.T) {
	d := db.CreateDB()
	a, _ := AddUser(d, "alice", "bar")

	Follow(d, a.UUID, a.UUID)

	alice, _ := d.FindUser(a.UUID)

	assert.Equal(t, alice.Following, []string{})
	assert.Equal(t, alice.Followers, []string{})
}

func TestAlreadyFollowingFail(t *testing.T) {
	d := db.CreateDB()
	a, _ := AddUser(d, "alice", "bar")
	b, _ := AddUser(d, "bill", "bar")

	Follow(d, a.UUID, b.UUID)

	alice, _ := d.FindUser(a.UUID)
	bill, _ := d.FindUser(b.UUID)

	assert.Equal(t, alice.Following, []string{b.UUID})
	assert.Equal(t, bill.Followers, []string{a.UUID})

	Follow(d, a.UUID, b.UUID)

	alice, _ = d.FindUser(a.UUID)
	bill, _ = d.FindUser(b.UUID)

	assert.Equal(t, alice.Following, []string{b.UUID})
	assert.Equal(t, bill.Followers, []string{a.UUID})
}

func TestNotFoundPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Follow should have panicked but did not")
		}
	}()

	d := db.CreateDB()
	a, _ := AddUser(d, "alice", "bar")
	b := NewUser("bill") // Not added to DB

	Follow(d, a.UUID, b.UUID)
}

func TestLoginSuccess(t *testing.T) {
	d := db.CreateDB()
	SeedUsers(d)
	u, err := Login(d, "alice", "bar")
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
	u, err := Login(d, "alice", "")
	assert.Nil(t, u)
	assert.Equal(t, err.Error(), "Name and password must be provided")
}

func TestLoginMissingUser(t *testing.T) {
	d := db.CreateDB()
	SeedUsers(d)
	u, err := Login(d, "nobody", "foo")
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
