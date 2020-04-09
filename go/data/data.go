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
	InsertTodo() error
	SelectAllTodos(name string) ([]types.TodoData, error)
	SelectByPriority(priority int, name string) error
	SelectNonPriority(name string) error
	SelectByCategory(category string, name string) error
	SelectByID(id int, name string) error
	DeleteByTitle(title int, name string) error
	DeleteByPriority(priority int, name string) error
	DeleteInactive(name string) error
	DeleteByID(id int, name string) error
	UpdateTitle(id int, newTitle string, name string)
	UpdatePriority(id int, newPriority int, name string)
	UpdateActive(id int, newActive bool, name string) error
}

//StoreType is the struct holding the db connection
type StoreType struct {
	DAO *sql.DB
}

//InsertTodo adds a brand new, fresh, shiny, little todo item to the todo list
func (store *StoreType) InsertTodo(todo types.TodoData) error {
	_, err := store.DAO.Exec(`
INSERT INTO Todos (acct_name, title, body, category, item_priority, publish_date, active) VALUES (?, ?, ?, ?, ?, ?, ?)`, todo.Name, todo.Title, todo.Body, todo.Category, todo.Priority, time.Now(), true)
	if err != nil {
		log.Errorf("Error inserting todo item: %v", err)
		return err
	}
	return nil
}

//SelectAllTodos selects all todo items from the db.  Providing a value of true for active will cause todo items to be returned only if they are actice
func (store *StoreType) SelectAllTodos(name string) ([]types.TodoData, error) {
	results, err := store.DAO.Query(`SELECT * FROM Todos WHERE acct_name = ?`, name)
	if err != nil {
		log.Errorf("Error querying mysql: %v", err)
		return nil, err
	}
	var tags []types.TodoData
	for results.Next() {
		var tag types.TodoData
		err = results.Scan(&tag.ID, &tag.Name, &tag.Title, &tag.Body, &tag.Category, &tag.Priority, &tag.PublishDate, &tag.Active)
		if err != nil {
			log.Warnf("Error selecting single row: %v", err)
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

//SelectActives selects all todo items from the db.  Providing a value of true for active will cause todo items to be returned only if they are actice
func (store *StoreType) SelectActives(active bool, name string) ([]types.TodoData, error) {
	results, err := store.DAO.Query(`SELECT * FROM Todos WHERE active = ? AND acct_name = ?`, active, name)
	if err != nil {
		log.Errorf("Error querying mysql: %v", err)
		return nil, err
	}
	var tags []types.TodoData
	for results.Next() {
		var tag types.TodoData
		err = results.Scan(&tag.ID, &tag.Name, &tag.Title, &tag.Body, &tag.Category, &tag.Priority, &tag.PublishDate, &tag.Active)
		if err != nil {
			log.Warnf("Error selecting single row: %v", err)
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

//SelectByPriority returns all todo items at or above the priority specified ...
func (store *StoreType) SelectByPriority(priority int, name string) ([]types.TodoData, error) {
	results, err := store.DAO.Query(`SELECT * FROM Todos WHERE item_priority <= ? AND acct_name = ?`, priority, name)
	if err != nil {
		log.Errorf("Error querying mysql: %v", err)
		return nil, err
	}
	var tags []types.TodoData
	for results.Next() {
		var tag types.TodoData
		err = results.Scan(&tag.ID, &tag.Name, &tag.Title, &tag.Body, &tag.Category, &tag.Priority, &tag.PublishDate, &tag.Active)
		if err != nil {
			log.Warnf("Error selecting single row: %v", err)
			return nil, err
		}
		if tag.Priority != 0 {
			tags = append(tags, tag)
		}
	}
	return tags, nil
}

//SelectNonPriority returns all todo items that do not have a priority specified (priority == 0)
func (store *StoreType) SelectNonPriority(name string) ([]types.TodoData, error) {
	results, err := store.DAO.Query(`SELECT * FROM Todos WHERE item_priority = 0 AND acct_name = ?`, name)
	if err != nil {
		log.Errorf("Error querying mysql: %v", err)
		return nil, err
	}
	var tags []types.TodoData
	for results.Next() {
		var tag types.TodoData
		err = results.Scan(&tag.ID, &tag.Name, &tag.Title, &tag.Body, &tag.Category, &tag.Priority, &tag.PublishDate, &tag.Active)
		if err != nil {
			log.Warnf("Error selecting single row: %v", err)
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

//SelectByCategory ...
func (store *StoreType) SelectByCategory(category string, name string) ([]types.TodoData, error) {
	results, err := store.DAO.Query(`SELECT * FROM Todos WHERE category = ? AND acct_name = ?`, category, name)
	if err != nil {
		log.Errorf("Error querying mysql: %v", err)
		return nil, err
	}
	var tags []types.TodoData
	for results.Next() {
		var tag types.TodoData
		err = results.Scan(&tag.ID, &tag.Name, &tag.Title, &tag.Body, &tag.Category, &tag.Priority, &tag.PublishDate, &tag.Active)
		if err != nil {
			log.Warnf("Error selecting single row: %v", err)
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

//SelectByID returns the todo item associated with the given id
func (store *StoreType) SelectByID(id int, name string) (types.TodoData, error) {
	var tag types.TodoData
	err := store.DAO.QueryRow(`SELECT * FROM Todos WHERE id = ? AND acct_name = ?`, id, name).Scan(&tag.ID, &tag.Name, &tag.Title, &tag.Body, &tag.Category, &tag.Priority, &tag.PublishDate, &tag.Active)
	if err != nil {
		log.Errorf("Error querying mysql: %v", err)
		return tag, err
	}
	return tag, nil
}

//DeleteByTitle deletes all todo items with the specified title
func (store *StoreType) DeleteByTitle(title string, name string) error {
	_, err := store.DAO.Query(`DELETE FROM Todos WHERE title = ? AND acct_name = ?`, title, name)
	return err
}

//DeleteByPriority deletes all todo items at the given priority level
func (store *StoreType) DeleteByPriority(priority int, name string) error {
	_, err := store.DAO.Query(`DELETE FROM Todos WHERE item_priority = ? AND acct_name = ?`, priority, name)
	return err
}

//DeleteInactive deletes all todo items at the given priority level
func (store *StoreType) DeleteInactive(name string) error {
	_, err := store.DAO.Query(`DELETE FROM Todos WHERE active = false AND acct_name = ?`, name)
	return err
}

//DeleteByID deletes a todo item that has the given ID
func (store *StoreType) DeleteByID(id int, name string) error {
	_, err := store.DAO.Query(`DELETE FROM Todos WHERE id = ? AND acct_name = ?`, id, name)
	return err
}

//UpdateTitle updates the title of a todo iten based on its id
func (store *StoreType) UpdateTitle(id int, newTitle string, name string) error {
	var count int
	err := store.DAO.QueryRow(`SELECT count(*) FROM Todos WHERE id = ? AND acct_name = ?`, id, name).Scan(&count)
	if count == 0 {
		return errors.New("ID does not exist")
	}

	_, err = store.DAO.Exec(`UPDATE Todos SET title = ? WHERE id = ? AND acct_name = ?`, newTitle, id, name)
	if err != nil {
		log.Errorf("Error updating title: %v", err)
		return err
	}
	return nil
}

//UpdatePriority updates the priority level of a todo iten based on its id
func (store *StoreType) UpdatePriority(id int, newPriority int, name string) error {
	var count int
	err := store.DAO.QueryRow(`SELECT count(*) FROM Todos WHERE id = ? AND acct_name = ?`, id, name).Scan(&count)
	if count == 0 {
		return errors.New("ID does not exist")
	}

	_, err = store.DAO.Exec(`UPDATE Todos SET item_priority = ? WHERE id = ? AND acct_name = ?`, newPriority, id, name)
	if err != nil {
		log.Errorf("Error updating rating: %v", err)
		return err
	}
	return nil
}

//UpdateActive updates whether a todo item is active or not, based on its id
func (store *StoreType) UpdateActive(id int, newActive bool, name string) error {
	var count int
	err := store.DAO.QueryRow(`SELECT count(*) FROM Todos WHERE id = ? AND acct_name = ?`, id, name).Scan(&count)
	if count == 0 {
		return errors.New("Title does not exist")
	}

	_, err = store.DAO.Exec(`UPDATE Todos SET active = ? WHERE id = ? AND acct_name = ?`, newActive, id, name)
	if err != nil {
		log.Errorf("Error updating rating: %v", err)
		return err
	}
	return nil
}
