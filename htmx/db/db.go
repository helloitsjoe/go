package db

type user struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

type User struct {
	Username string
}

// TODO: not string
func CreateDB() string {
	// schema := &memdb.DBSchema{
	// 	Tables: map[string]*memdb.TableSchema{
	// 		"user": &memdb.TableSchema{
	// 			Name: "user",
	// 			Indexes: map[string]*memdb.IndexSchema{
	// 				"username": &memdb.IndexSchema{
	// 					Name:    "username",
	// 					Unique:  true,
	// 					Indexer: &memdb.StringFieldIndex{Field: "Username"},
	// 				},
	// 			},
	// 		},
	// 	},
	// }

	// _, err := memdb.NewMemDB(schema)
	// if err != nil {
	// 	panic(err)
	// }

	return ""
}
