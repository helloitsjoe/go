package user

import (
	"htmx/db"
	"testing"
)

func TestLogin(t *testing.T) {
	// TODO: Mock DB
	mockDb := db.CreateDB()
	u, err := Login(mockDb, "Alice", "bar")

}
