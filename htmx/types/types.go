package types

import "github.com/google/uuid"

// type user struct {
// 	Username string `form:"username"`
// 	Password string `form:"password"`
// }

type User struct {
	Username  string
	Followers []uuid.UUID
	Following []uuid.UUID
	// LikedPosts []uuid.UUID
	// Posts      []Post
	UUID uuid.UUID
}

// type Post struct {
// 	Content    string
// 	Author     uuid.UUID
// 	Likes      []uuid.UUID
// 	Responses  []uuid.UUID
// 	ResponseTo uuid.UUID
// 	UUID       uuid.UUID
// }
