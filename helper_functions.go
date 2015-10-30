package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/ottermad/odotgo/todo"
)

func NewDB() *sql.DB {
	db, err := sql.Open("sqlite3", "db.sqlite3")
	if err != nil {
		panic(err)
	}

	todo.CreateTables(db)

	return db
}
