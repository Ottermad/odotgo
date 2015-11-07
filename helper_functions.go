package main

import (
	"database/sql"
	"fmt"
	_ "github.com/ottermad/odotgo/Godeps/_workspace/src/github.com/lib/pq"
	"github.com/ottermad/odotgo/todo"
	"os"
)

func NewDB() *sql.DB {
	localDev := os.Getenv("LOCAL_DEV")

	var db *sql.DB
	var err error
	fmt.Fprint(os.Stdout, "STUFF")
	if localDev == "TRUE" {
		db, err = sql.Open("postgres", "dbname=ODOT sslmode=disable")
		if err != nil {
			panic(err)
		}
	} else {
		url := os.Getenv("DATABASE_URL")
		fmt.Fprint(os.Stdout, url)
		db, err = sql.Open("postgres", url)
		if err != nil {
			panic(err)
		}
	}

	todo.CreateTables(db)

	return db
}
