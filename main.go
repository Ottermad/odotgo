package odot

import (
	"net/http"
	"odot/todo"

	"github.com/julienschmidt/httprouter"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	r := httprouter.New()

	r.GET("/", todo.TodoListListing)

	r.GET("/todo", todo.TodoListListing)
	r.GET("/todo/view/:id", todo.TodoListDetail)
	r.GET("/todo/add", todo.AddTodoListGet)
	r.POST("/todo/add", todo.AddTodoListPost)

	http.ListenAndServe(":8080", r)
}
