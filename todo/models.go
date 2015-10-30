package todo

import (
	"database/sql"
	"fmt"
)

// Helper Functions
func CreateTables(db *sql.DB) {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS TBL_TODO (ID INTEGER PRIMARY KEY AUTOINCREMENT, TITLE VARCHAR(100) NOT NULL UNIQUE, DESCRIPTION VARCHAR(250) NOT NULL)")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS TBL_TODO_ITEM (ID INTEGER PRIMARY KEY AUTOINCREMENT, TODO_LIST INTEGER NOT NULL, CONTENT VARCHAR(250) NOT NULL, FOREIGN KEY (TODO_LIST) REFERENCES TBL_TODO(ID), UNIQUE (TODO_LIST, CONTENT))")

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
	_, err := db.Exec("INSERT INTO TBL_TODO (TITLE, DESCRIPTION) VALUES (?, ?)", title, description)

	if err != nil {
		return err
	}

	// Fetch row and set id
	err = db.QueryRow("SELECT ID FROM TBL_TODO WHERE TITLE = ?", title).Scan(&todoList.Id)

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

func (todoList *TodoList) FindById(id int, db *sql.DB) error {
	err := db.QueryRow("SELECT ID, TITLE, DESCRIPTION FROM TBL_TODO WHERE ID = ?", id).Scan(
		&todoList.Id,
		&todoList.Title,
		&todoList.Description)
	if err != nil {
		return err
	}

	rows, err := db.Query("SELECT ID, CONTENT FROM TBL_TODO_ITEM WHERE TODO_LIST = ?", todoList.Id)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer rows.Close()

	for rows.Next() {
		listItem := TodoListItem{}
		if err = rows.Scan(&listItem.Id, &listItem.Content); err != nil {
			return err
		}
		fmt.Println(listItem)
		todoList.Items = append(todoList.Items, listItem)
	}
	todoList.db = db

	return nil
}

func (todoList *TodoList) FindByTitle(title string, db *sql.DB) error {
	err := db.QueryRow("SELECT ID, TITLE, DESCRIPTION FROM TBL_TODO WHERE TITLE = ?", title).Scan(
		&todoList.Id,
		&todoList.Title,
		&todoList.Description)
	if err != nil {
		return err
	}

	rows, err := db.Query("SELECT ID, CONTENT FROM TBL_TODO_ITEM WHERE TODO_LIST = ?", todoList.Id)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer rows.Close()

	for rows.Next() {
		listItem := TodoListItem{}
		if err = rows.Scan(&listItem.Id, &listItem.Content); err != nil {
			return err
		}
		fmt.Println(listItem)
		todoList.Items = append(todoList.Items, listItem)
	}

	todoList.db = db

	return nil
}

func (todoList *TodoList) AddItem(content string) error {
	// Insert into database

	_, err := todoList.db.Exec("INSERT INTO TBL_TODO_ITEM (TODO_LIST, CONTENT) VALUES (?, ?)", todoList.Id, content)
	if err != nil {
		fmt.Println("INSERTING ERROR")
		return err
	}

	// Get id
	var id int
	err = todoList.db.QueryRow("SELECT ID FROM TBL_TODO_ITEM WHERE TODO_LIST = ? AND CONTENT = ?", todoList.Id, content).Scan(&id)

	if err != nil {
		fmt.Println("SELECT ERROR")
		return err
	}

	// crate struct and set attributes
	todoItem := TodoListItem{Id: id, Content: content}
	todoList.Items = append(todoList.Items, todoItem)

	return nil
}
