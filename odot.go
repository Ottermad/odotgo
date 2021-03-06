package main

import (
	"net/http"
	"os"

	"github.com/ottermad/odotgo/todo"

	"github.com/ottermad/odotgo/Godeps/_workspace/src/github.com/julienschmidt/httprouter"
)

func main() {
	db := NewDB()
	r := httprouter.New()

	r.ServeFiles("/static/*filepath", http.Dir("public"))

	r.GET("/", todo.TodoListListing(db))
	r.GET("/todo", todo.TodoListListing(db))

	r.GET("/todo/add", todo.AddTodoListGet(db))
	r.POST("/todo/add", todo.AddTodoListPost(db))

	r.GET("/todo/view/:id", todo.TodoListDetail(db))

	r.GET("/todo/edit/:id", todo.EditTodoListGet(db))
	r.POST("/todo/edit/:id", todo.EditTodoListPost(db))

	r.GET("/todo/add-item/:id", todo.AddTodoListItemGet(db))
	r.POST("/todo/add-item/:id", todo.AddTodoListItemPost(db))

	r.GET("/todo/delete/:id", todo.DeleteTodoList(db))

	r.GET("/todo/delete-item/:todolistid/:todoitemid", todo.DeleteTodoListItem(db))

	if os.Getenv("LOCAL_DEV") == "TRUE" {
		http.ListenAndServe(":8080", r)
	} else {
		http.ListenAndServe(":"+os.Getenv("PORT"), r)
	}

}
