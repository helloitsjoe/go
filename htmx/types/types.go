package types

// type user struct {
// 	Username string `form:"username"`
// 	Password string `form:"password"`
// }

// UUIDs need to be converted to strings for memdb,
// otherwise schema gives a length mismatch error
type User struct {
	Username  string
	Followers []string
	Following []string
	// LikedPosts []uuid.UUID
	// Posts      []Post
	UUID string
}

// type Post struct {
// 	Content    string
// 	Author     uuid.UUID
// 	Likes      []uuid.UUID
// 	Responses  []uuid.UUID
// 	ResponseTo uuid.UUID
// 	UUID       uuid.UUID
// }
