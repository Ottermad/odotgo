package todo

import (
	"database/sql"
	"errors"
)

// Helper Functions
func CreateTables(db *sql.DB) {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS "TBL_TODO" (
			ID SERIAL PRIMARY KEY, 
			TITLE VARCHAR(100) NOT NULL UNIQUE, 
			DESCRIPTION VARCHAR(250) NOT NULL
		)
	`)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS "TBL_TODO_ITEM" (
			ID SERIAL PRIMARY KEY, 
			TODO_LIST INTEGER NOT NULL, 
			CONTENT VARCHAR(250) NOT NULL, 
			FOREIGN KEY (TODO_LIST) REFERENCES "TBL_TODO"(ID) ON DELETE CASCADE, 
			UNIQUE (TODO_LIST, CONTENT)
		)
	`)

	if err != nil {
		panic(err)
	}
}

// Structs
type TodoList struct {
	Id          int
	Title       string
	Description string
	Items       []TodoListItem
	db          *sql.DB
}

type TodoListItem struct {
	Id      int
	Content string
}

// Methods
func (todoList *TodoList) New(title string, description string, db *sql.DB) error {
	// Insert into TBL_TODO
	_, err := db.Exec(`INSERT INTO "TBL_TODO" (TITLE, DESCRIPTION) VALUES ($1, $2)`, title, description)

	if err != nil {
		return err
	}

	// Fetch row and set id
	err = db.QueryRow(`SELECT ID FROM "TBL_TODO" WHERE TITLE = $1`, title).Scan(&todoList.Id)

	if err != nil {
		return err
	}

	// Set other attributes
	todoList.Title = title
	todoList.Description = description
	todoList.db = db
	todoList.Items = []TodoListItem{}

	return nil
}

func (todoList *TodoList) Delete() error {
	_, err := todoList.db.Exec(`DELETE FROM "TBL_TODO" WHERE ID = $1"`, todoList.Id)
	if err != nil {
		return err
	}

	todoList.Id = -1
	todoList.Title = ""
	todoList.Description = ""
	todoList.Items = nil
	todoList.db = nil

	return nil
}

func (todoList *TodoList) FindById(id int, db *sql.DB) error {
	err := db.QueryRow(`SELECT ID, TITLE, DESCRIPTION FROM "TBL_TODO" WHERE ID = $1`, id).Scan(
		&todoList.Id,
		&todoList.Title,
		&todoList.Description)
	if err != nil {
		return err
	}

	rows, err := db.Query(`SELECT ID, CONTENT FROM "TBL_TODO_ITEM" WHERE TODO_LIST = $1`, todoList.Id)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		listItem := TodoListItem{}
		if err = rows.Scan(&listItem.Id, &listItem.Content); err != nil {
			return err
		}
		todoList.Items = append(todoList.Items, listItem)
	}
	todoList.db = db

	return nil
}

func (todoList *TodoList) FindByTitle(title string, db *sql.DB) error {
	err := db.QueryRow(`SELECT ID, TITLE, DESCRIPTION FROM "TBL_TODO" WHERE TITLE = $1`, title).Scan(
		&todoList.Id,
		&todoList.Title,
		&todoList.Description)
	if err != nil {
		return err
	}

	rows, err := db.Query(`SELECT ID, CONTENT FROM "TBL_TODO_ITEM" WHERE TODO_LIST = $1`, todoList.Id)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		listItem := TodoListItem{}
		if err = rows.Scan(&listItem.Id, &listItem.Content); err != nil {
			return err
		}
		todoList.Items = append(todoList.Items, listItem)
	}

	todoList.db = db

	return nil
}

func (todoList *TodoList) AddItem(content string) error {
	// Insert into database

	_, err := todoList.db.Exec(`INSERT INTO "TBL_TODO_ITEM" (TODO_LIST, CONTENT) VALUES ($1, $2)`, todoList.Id, content)
	if err != nil {
		return err
	}

	// Get id
	var id int
	err = todoList.db.QueryRow(`SELECT ID FROM "TBL_TODO_ITEM" WHERE TODO_LIST = $1 AND CONTENT = $2`, todoList.Id, content).Scan(&id)

	if err != nil {
		return err
	}

	// crate struct and set attributes
	todoItem := TodoListItem{Id: id, Content: content}
	todoList.Items = append(todoList.Items, todoItem)

	return nil
}

func (todoList *TodoList) DeleteItem(id int) error {
	// Check id is valid
	is_valid := false
	var item_index int
	for index, item := range todoList.Items {
		if item.Id == id {
			is_valid = true
			item_index = index
		}
	}

	if !is_valid {
		return errors.New("id not in items of TodoList")
	}

	_, err := todoList.db.Exec(`DELETE FROM "TBL_TODO_ITEM" WHERE ID = $1`, id)

	if err != nil {
		return err
	}

	// Remove from items
	todoList.Items = append(todoList.Items[:item_index], todoList.Items[item_index+1:]...)

	return nil
}

func (todoList *TodoList) UpdateTitle(title string) error {
	_, err := todoList.db.Exec(`UPDATE "TBL_TODO" SET TITLE = $1 WHERE ID = $2`, title, todoList.Id)

	if err != nil {
		return err
	}

	todoList.Title = title

	return nil
}

func (todoList *TodoList) UpdateDescription(description string) error {
	_, err := todoList.db.Exec(`UPDATE "TBL_TODO" SET DESCRIPTION = $1 WHERE ID = $2`, description, todoList.Id)

	if err != nil {
		return err
	}

	todoList.Description = description

	return nil
}
