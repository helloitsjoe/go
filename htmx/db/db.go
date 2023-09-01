package db

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func CreateDB() {
	db, err := sql.Open("sqlite", dsnURI)
}
