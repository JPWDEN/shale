package data

import (
	"database/sql"
	"errors"
	"time"

	"github.com/bdlm/log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/shale/go/types"
)

//Store is the interface defining the object for db functions
type Store interface {
	SQLTest(todo types.TodoData) error
	InsertTodo() error
	SelectAllTodos() ([]types.TodoData, error)
	SelectByPriority(priority int) error
	SelectByCategory(category string) error
	SelectByID(id int) error
	DeleteByTitle(title int) error
	DeleteByPriority(priority int) error
	DeleteByID(id int) error
	UpdateTitle(id int, newTitle string)
	UpdatePriority(id int, newPriority int)
	UpdateActive(id int, newActive bool) error
}

//StoreType is the struct holding the db connection
type StoreType struct {
	DAO *sql.DB
}

//SQLTest tests the db connection
func (store *StoreType) SQLTest(todo types.TodoData) error {
	_, err := store.DAO.Exec(`
INSERT INTO todos VALUES (?, ?, ?, ?, ?, ?)`, todo.ID)
	if err != nil {
		log.Error(err)
	}
	return nil
}

//InsertTodo adds a brand new, fresh, shiny, little todo item to the todo list
func (store *StoreType) InsertTodo(todo types.TodoData) error {
	_, err := store.DAO.Exec(`
INSERT INTO Todos (title, body, category, item_priority, publish_date, active) VALUES (?, ?, ?, ?, ?, ?)`, todo.Title, todo.Body, todo.Category, todo.Priority, time.Now(), true)
	if err != nil {
		log.Errorf("Error inserting todo item: %v", err)
		return err
	}
	return nil
}

//SelectAllTodos selects all todo items from the db.  Providing a value of true for active will cause todo items to be returned only if they are actice
func (store *StoreType) SelectAllTodos(active bool) ([]types.TodoData, error) {
	var results *sql.Rows
	var err error
	if active {
		results, err = store.DAO.Query(`SELECT * FROM Todos WHERE active = TRUE`)
	} else {
		results, err = store.DAO.Query(`SELECT * FROM Todos`)
	}
	if err != nil {
		log.Errorf("Error querying mysql: %v", err)
		return nil, err
	}
	var tags []types.TodoData
	for results.Next() {
		var tag types.TodoData
		err = results.Scan(&tag.ID, &tag.Title, &tag.Body, &tag.Category, &tag.Priority, &tag.PublishDate, &tag.Active)
		if err != nil {
			log.Warnf("Error selecting single row: %v", err)
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

//SelectByPriority returns all todo items at or above the priority specified ...
func (store *StoreType) SelectByPriority(priority int) ([]types.TodoData, error) {
	results, err := store.DAO.Query(`SELECT * FROM Todos WHERE item_priority <= ?`, priority)
	if err != nil {
		log.Errorf("Error querying mysql: %v", err)
		return nil, err
	}
	var tags []types.TodoData
	for results.Next() {
		var tag types.TodoData
		err = results.Scan(&tag.ID, &tag.Title, &tag.Body, &tag.Category, &tag.Priority, &tag.PublishDate, &tag.Active)
		if err != nil {
			log.Warnf("Error selecting single row: %v", err)
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

//SelectByCategory ...
func (store *StoreType) SelectByCategory(category string) ([]types.TodoData, error) {
	results, err := store.DAO.Query(`SELECT * FROM Todos WHERE category = ?`, category)
	if err != nil {
		log.Errorf("Error querying mysql: %v", err)
		return nil, err
	}
	var tags []types.TodoData
	for results.Next() {
		var tag types.TodoData
		err = results.Scan(&tag.ID, &tag.Title, &tag.Body, &tag.Category, &tag.Priority, &tag.PublishDate, &tag.Active)
		if err != nil {
			log.Warnf("Error selecting single row: %v", err)
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

//SelectByID returns the todo item associated with the given id
func (store *StoreType) SelectByID(id int) (types.TodoData, error) {
	var tag types.TodoData
	result, err := store.DAO.Query(`SELECT * FROM Todos WHERE id = ?`, id)
	if err != nil {
		log.Errorf("Error querying mysql: %v", err)
		return tag, err
	}
	err = result.Scan(&tag.ID, &tag.Title, &tag.Body, &tag.Category, &tag.Priority, &tag.PublishDate, &tag.Active)
	if err != nil {
		log.Warnf("Error selecting single row: %v", err)
		return tag, err
	}
	return tag, nil
}

//DeleteByPriority deletes all todo items at the given priority level
func (store *StoreType) DeleteByPriority(priority int) error {
	_, err := store.DAO.Query(`DELETE FROM Todos WHERE PRIORITY = ?`, priority)
	return err
}

//DeleteByTitle deletes all todo items with the specified title
func (store *StoreType) DeleteByTitle(title string) error {
	_, err := store.DAO.Query(`DELETE FROM Todos WHERE title = ?`, title)
	return err
}

//DeleteByID deletes a todo item that has the given ID
func (store *StoreType) DeleteByID(id int) error {
	_, err := store.DAO.Query(`DELETE FROM Todos WHERE id = ?`, id)
	return err
}

//UpdateTitle updates the title of a todo iten based on its id
func (store *StoreType) UpdateTitle(id int, newTitle string) error {
	var count int
	err := store.DAO.QueryRow(`SELECT count(*) FROM todos WHERE id = ?`, id).Scan(&count)
	if count == 0 {
		return errors.New("Title does not exist")
	}

	_, err = store.DAO.Exec(`UPDATE todos SET title = ? WHERE id = ?`, newTitle, id)
	if err != nil {
		log.Errorf("Error updating rating: %v", err)
		return err
	}
	return nil
}

//UpdatePriority updates the priority level of a todo iten based on its id
func (store *StoreType) UpdatePriority(id int, newPriority int) error {
	var count int
	err := store.DAO.QueryRow(`SELECT count(*) FROM todos WHERE id = ?`, id).Scan(&count)
	if count == 0 {
		return errors.New("Title does not exist")
	}

	_, err = store.DAO.Exec(`UPDATE todos SET priority = ? WHERE id = ?`, newPriority, id)
	if err != nil {
		log.Errorf("Error updating rating: %v", err)
		return err
	}
	return nil
}

//UpdateActive updates whether a todo item is active or not, based on its id
func (store *StoreType) UpdateActive(id int, newActive bool) error {
	var count int
	err := store.DAO.QueryRow(`SELECT count(*) FROM todos WHERE id = ?`, id).Scan(&count)
	if count == 0 {
		return errors.New("Title does not exist")
	}

	_, err = store.DAO.Exec(`UPDATE todos SET active = ? WHERE id = ?`, newActive, id)
	if err != nil {
		log.Errorf("Error updating rating: %v", err)
		return err
	}
	return nil
}
