package todo

import (
	"database/sql"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"net/http"
	"path"
	"strconv"
)

func AddTodoListPost(db *sql.DB) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		todoList := TodoList{}
		err := todoList.New(r.FormValue("name"), r.FormValue("description"), db)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(rw, r, "/todo", http.StatusOK)
	}
}

func AddTodoListGet(db *sql.DB) httprouter.Handle {

	return func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		fp := path.Join("public", "todo_add.html")
		tmpl, err := template.ParseFiles(fp)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := tmpl.Execute(rw, nil); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}

	}

}
func TodoListListing(db *sql.DB) httprouter.Handle {

	return func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

		// Get all todo_lists
		rows, err := db.Query("SELECT ID, TITLE FROM TBL_TODO")
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var items []TodoList
		for rows.Next() {
			list := TodoList{}
			if err = rows.Scan(&list.Id, &list.Title); err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}
			items = append(items, list)
		}

		fp := path.Join("public", "todo_listing.html")
		tmpl, err := template.ParseFiles(fp)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := tmpl.Execute(rw, items); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}

	}

}

func TodoListDetail(db *sql.DB) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		groceryList := TodoList{}
		id, err := strconv.Atoi(p.ByName("id"))
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		err = groceryList.FindById(id, db)
		fmt.Println(err)

		fp := path.Join("public", "todo_detail.html")
		tmpl, err := template.ParseFiles(fp)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := tmpl.Execute(rw, groceryList); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	}
}
