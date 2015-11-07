package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/ottermad/odotgo/todo"

	"github.com/ottermad/odotgo/Godeps/_workspace/src/github.com/julienschmidt/httprouter"
)

func main() {
	db := NewDB()
	fmt.Fprintln(os.Stdout, "Run NewDB")
	r := httprouter.New()
	fmt.Fprintln(os.Stdout, "Create httprouter")

	r.ServeFiles("/static/*filepath", http.Dir("public"))
	fmt.Fprintln(os.Stdout, "Serving Status")

	r.GET("/", todo.TodoListListing(db))
	r.GET("/todo", todo.TodoListListing(db))
	fmt.Fprintln(os.Stdout, "TodoListListing")

	r.GET("/todo/add", todo.AddTodoListGet(db))
	fmt.Fprintln(os.Stdout, "AddTodoListGet")

	r.POST("/todo/add", todo.AddTodoListPost(db))
	fmt.Fprintln(os.Stdout, "AddTodoListPost")

	r.GET("/todo/view/:id", todo.TodoListDetail(db))
	fmt.Fprintln(os.Stdout, "TodoListDetail")

	r.GET("/todo/edit/:id", todo.EditTodoListGet(db))
	fmt.Fprintln(os.Stdout, "EditTodoListGet")

	r.POST("/todo/edit/:id", todo.EditTodoListPost(db))
	fmt.Fprintln(os.Stdout, "EditTodoListPost")

	r.GET("/todo/add-item/:id", todo.AddTodoListItemGet(db))
	fmt.Fprintln(os.Stdout, "AddTodoListItemGet")

	r.POST("/todo/add-item/:id", todo.AddTodoListItemPost(db))
	fmt.Fprintln(os.Stdout, "AddTodoListItemPost")

	r.GET("/todo/delete/:id", todo.DeleteTodoList(db))
	fmt.Fprintln(os.Stdout, "DeleteTodoList")

	r.GET("/todo/delete-item/:todolistid/:todoitemid", todo.DeleteTodoListItem(db))
	fmt.Fprintln(os.Stdout, "DeleteTodoListItem")
	fmt.Fprintln(os.Stdout, os.Getenv("PORT"))

	if os.Getenv("LOCAL_DEV") == "TRUE" {
		http.ListenAndServe(":8080", r)
	} else {
		http.ListenAndServe(":"+os.Getenv("PORT"), r)
		fmt.Fprintln(os.Stdout, "Running")

	}

}
