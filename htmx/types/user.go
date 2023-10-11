package types

import "github.com/google/uuid"

// type user struct {
// 	Username string `form:"username"`
// 	Password string `form:"password"`
// }

type User struct {
	Username string
	UUID     uuid.UUID
}
