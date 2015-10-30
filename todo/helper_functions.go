package todo

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func NewDB() *sql.DB {
	db, err := sql.Open("sqlite3", "db.sqlite3")
	if err != nil {
		panic(err)
	}

	CreateTables(db)

	return db
}
