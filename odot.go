package main

import (
	"net/http"

	"github.com/ottermad/odotgo/todo"

	"github.com/julienschmidt/httprouter"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db := NewDB()
	r := httprouter.New()

	r.GET("/", todo.TodoListListing(db))

	r.GET("/todo", todo.TodoListListing(db))
	r.GET("/todo/view/:id", todo.TodoListDetail(db))
	r.GET("/todo/add", todo.AddTodoListGet(db))
	r.POST("/todo/add", todo.AddTodoListPost(db))

	http.ListenAndServe(":8080", r)
}
