package todo

import (
	"database/sql"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"net/http"
	"path"
	"strconv"
)

const TEMPLATE_DIR string = "templates"

func AddTodoListPost(db *sql.DB) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		todoList := TodoList{}
		err := todoList.New(r.FormValue("title"), r.FormValue("description"), db)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(rw, r, "/todo", 301)
	}
}

func AddTodoListGet(db *sql.DB) httprouter.Handle {

	return func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		fp := path.Join(TEMPLATE_DIR, "todo_add.html")
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
		rows, err := db.Query(`SELECT ID, TITLE FROM "TBL_TODO"`)
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

		fp := path.Join(TEMPLATE_DIR, "todo_listing.html")
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

		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		fp := path.Join(TEMPLATE_DIR, "todo_detail.html")
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

func EditTodoListGet(db *sql.DB) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		todoList := TodoList{}
		id, err := strconv.Atoi(p.ByName("id"))
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		err = todoList.FindById(id, db)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusNotFound)
			return
		}

		fp := path.Join(TEMPLATE_DIR, "todo_edit.html")
		tmpl, err := template.ParseFiles(fp)

		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := tmpl.Execute(rw, todoList); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	}
}

func EditTodoListPost(db *sql.DB) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		todoList := TodoList{}
		id, err := strconv.Atoi(p.ByName("id"))
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		err = todoList.FindById(id, db)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusNotFound)
			return
		}

		if todoList.Title != r.FormValue("title") {
			err = todoList.UpdateTitle(r.FormValue("title"))
			if err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		if todoList.Description != r.FormValue("description") {
			err = todoList.UpdateDescription(r.FormValue("description"))
			if err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		http.Redirect(rw, r, "/todo/view/"+strconv.Itoa(todoList.Id), 301)
	}
}

func DeleteTodoList(db *sql.DB) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		todoList := TodoList{}
		id, err := strconv.Atoi(p.ByName("id"))
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}

		err = todoList.FindById(id, db)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}

		err = todoList.Delete()
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}

		http.Redirect(rw, r, "/todo", 301)
	}
}

func AddTodoListItemGet(db *sql.DB) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		todoList := TodoList{}
		id, err := strconv.Atoi(p.ByName("id"))
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}

		err = todoList.FindById(id, db)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusNotFound)
		}

		fp := path.Join(TEMPLATE_DIR, "todo_list_item_add.html")
		tmpl, err := template.ParseFiles(fp)

		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := tmpl.Execute(rw, todoList); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	}
}

func AddTodoListItemPost(db *sql.DB) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		todoList := TodoList{}
		id, err := strconv.Atoi(p.ByName("id"))
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}

		err = todoList.FindById(id, db)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusNotFound)
		}

		err = todoList.AddItem(r.FormValue("content"))
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}

		http.Redirect(rw, r, "/todo/view/"+strconv.Itoa(todoList.Id), 301)

	}
}

func DeleteTodoListItem(db *sql.DB) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		// Get todo list
		todoList := TodoList{}
		todoListId, err := strconv.Atoi(p.ByName("todolistid"))
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}

		err = todoList.FindById(todoListId, db)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusNotFound)
		}

		todoItemId, err := strconv.Atoi(p.ByName("todoitemid"))
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}

		err = todoList.DeleteItem(todoItemId)

		if err != nil {
			http.Error(rw, err.Error(), http.StatusNotFound)
		}

		http.Redirect(rw, r, "/todo/view/"+strconv.Itoa(todoList.Id), 301)
	}
}
