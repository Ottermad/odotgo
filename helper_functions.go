package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/ottermad/odotgo/todo"
	"os"
)

func NewDB() *sql.DB {
	localDev := os.Getenv("LOCAL_DEV")

	var db *sql.DB
	var err error

	if localDev == "TRUE" {
		db, err = sql.Open("postgres", "dbname=ODOT sslmode=disable")
		if err != nil {
			panic(err)
		}
	} else {
		url := os.Getenv("DATABASE_URL")
		db, err = sql.Open("postgres", url)
		if err != nil {
			panic(err)
		}
	}

	todo.CreateTables(db)

	return db
}
