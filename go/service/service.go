package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/bdlm/log"
	"github.com/shale/go/data"
	"github.com/shale/go/types"
)

//Server is the interface that defines the CRUD API object
type Server interface {
	HandleTodos(resp http.ResponseWriter, req *http.Request)
	AddTodo() error
	GetTodos(resp http.ResponseWriter, req *http.Request) error
	GetTodosByPriority(priority int, resp http.ResponseWriter, req *http.Request) error
	GetTodosByCategory(category string, resp http.ResponseWriter, req *http.Request) error
	GetTodosByID(id int, resp http.ResponseWriter, req *http.Request) error
	RemoveByTitle(resp http.ResponseWriter, req *http.Request) error
	RemoveByPriority(resp http.ResponseWriter, req *http.Request) error
	RemoveByID(resp http.ResponseWriter, req *http.Request) error
	ChangeTitle(id int, resp http.ResponseWriter, req *http.Request) error
	ChangePriority(id int, resp http.ResponseWriter, req *http.Request) error
	ChangeActive(id int, resp http.ResponseWriter, req *http.Request) error
}

//ServerType is the server object with db connection
type ServerType struct {
	DAO *data.StoreType
}

func encodeBody(resp http.ResponseWriter, req *http.Request, data interface{}) error {
	return json.NewEncoder(resp).Encode(data)
}

func decodeBody(req *http.Request, data interface{}) error {
	defer req.Body.Close()
	return json.NewDecoder(req.Body).Decode(data)
}

func respond(resp http.ResponseWriter, req *http.Request, status int, data interface{}) {
	resp.WriteHeader(status)
	if data != nil {
		encodeBody(resp, req, data)
	}
}

func respondErr(resp http.ResponseWriter, req *http.Request, status int, args ...interface{}) {
	respond(resp, req, status, map[string]interface{}{
		"error": map[string]interface{}{"message": fmt.Sprint(args...)},
	})
}

func respondHTTPErr(resp http.ResponseWriter, req *http.Request, status int) {
	respondErr(resp, req, status, http.StatusText(status))
}

//HandleTodos routes various API requests to the proper function
func (svr *ServerType) HandleTodos(resp http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	pathArgs := strings.Split(strings.Trim(path, "/"), "/")
	numArgs := len(pathArgs)
	queryMap, _ := url.ParseQuery(req.URL.RawQuery)
	fmt.Printf("PATHARGS %+v/%+v\n", pathArgs, queryMap)
	switch req.Method {
	case "GET":
		var err error
		switch pathArgs[1] {
		case "todos":
			if numArgs > 2 {
				active, err := strconv.ParseBool(pathArgs[2])
				if err == nil {
					err = svr.GetTodos(active, resp, req)
				}
			} else {
				err = svr.GetTodos(false, resp, req)
			}
		case "highs":
			if numArgs > 2 {
				priority, err := strconv.Atoi(pathArgs[2])
				if err == nil {
					err = svr.GetTodosByPriority(priority, resp, req)
				}
			} else {
				err = errors.New("Priority argument missing or malformed")
			}
		case "cat":
			if numArgs > 2 {
				err = svr.GetTodosByCategory(pathArgs[2], resp, req)
			}
		case "id":
			if numArgs > 2 {
				id, err := strconv.Atoi(pathArgs[2])
				if err == nil {
					err = svr.GetTodosByID(id, resp, req)
				}
			}
		}
		if err != nil {
			log.Errorf("GET Error: %v", err)
		}
		return
	case "POST":
		fmt.Printf("IN POST\n")
		var err error
		switch pathArgs[1] {
		case "add":
			err = svr.AddTodo(resp, req)
		case "title":
			if numArgs > 2 {
				id, err := strconv.Atoi(pathArgs[2])
				if err == nil {
					err = svr.ChangeTitle(id, resp, req)
				}
			} else {
				err = errors.New("Priority argument missing or malformed")
			}
		case "pri":
			if numArgs > 2 {
				id, err := strconv.Atoi(pathArgs[2])
				if err == nil {
					err = svr.ChangePriority(id, resp, req)
				}
			} else {
				err = errors.New("Priority argument missing or malformed")
			}
		case "active":
			if numArgs > 2 {
				id, err := strconv.Atoi(pathArgs[2])
				if err == nil {
					err = svr.ChangeActive(id, resp, req)
				}
			} else {
				err = errors.New("Priority argument missing or malformed")
			}
		}
		if err != nil {
			log.Errorf("POST Error: %v", err)
		}
		return
	case "DELETE":
		var err error
		switch pathArgs[1] {
		case "title":
			err = svr.RemoveByTitle(resp, req)
		case "pri":
			err = svr.RemoveByPriority(resp, req)
		case "active":
			err = svr.RemoveByID(resp, req)
		}
		if err != nil {
			log.Errorf("DELETE Error: %v", err)
		}
		return
	default:
		respondHTTPErr(resp, req, http.StatusBadRequest)
		return
	}

}

//AddTodo adds a new to do item to the to do list
func (svr *ServerType) AddTodo(resp http.ResponseWriter, req *http.Request) error {
	fmt.Printf("IN ADDTODO\n")
	var todo types.TodoData
	err := decodeBody(req, &todo)
	if err != nil {
		respondErr(resp, req, http.StatusBadRequest, " Failed to decode body: ", err)
		return err
	}
	err = svr.DAO.InsertTodo(todo)
	if err != nil {
		respondErr(resp, req, http.StatusInternalServerError, " Database error: ", err)
		return err
	}
	respond(resp, req, http.StatusOK, &types.ListStatus{
		Status: "Success",
		Info:   fmt.Sprintf("Title '%s' added", todo.Title),
	})
	return nil
}

//GetTodos returns all of the todo items from the list
func (svr *ServerType) GetTodos(active bool, resp http.ResponseWriter, req *http.Request) error {
	result, err := svr.DAO.SelectAllTodos(active)
	if err != nil {
		respondErr(resp, req, http.StatusInternalServerError, " Database error: ", err)
		return err
	}
	if result == nil {
		respond(resp, req, http.StatusNoContent, &result)
	} else {
		respond(resp, req, http.StatusOK, &result)
	}
	return nil
}

//GetTodosByPriority returns all todo items that have a priority higher (lower number) or equal to the one provided
func (svr *ServerType) GetTodosByPriority(priority int, resp http.ResponseWriter, req *http.Request) error {
	result, err := svr.DAO.SelectByPriority(priority)
	if err != nil {
		respondErr(resp, req, http.StatusInternalServerError, " Database error: ", err)
		return err
	}
	if result == nil {
		respond(resp, req, http.StatusNoContent, &result)
	} else {
		respond(resp, req, http.StatusOK, &result)
	}
	return nil
}

//GetTodosByCategory returns all todo items that exactly match the category provided
func (svr *ServerType) GetTodosByCategory(category string, resp http.ResponseWriter, req *http.Request) error {
	result, err := svr.DAO.SelectByCategory(category)
	if err != nil {
		respondErr(resp, req, http.StatusInternalServerError, " Database error: ", err)
		return err
	}
	if result == nil {
		respond(resp, req, http.StatusNoContent, &result)
	} else {
		respond(resp, req, http.StatusOK, &result)
	}
	return nil
}

//GetTodosByID returns the todo item associated with the given db id
func (svr *ServerType) GetTodosByID(id int, resp http.ResponseWriter, req *http.Request) error {
	result, err := svr.DAO.SelectByID(id)
	if err != nil {
		respondErr(resp, req, http.StatusInternalServerError, " Database error: ", err)
		return err
	}
	var empty types.TodoData
	if result == empty {
		respond(resp, req, http.StatusNoContent, &result)
	} else {
		respond(resp, req, http.StatusOK, &result)
	}
	return nil
}

//RemoveByTitle removes all todo items with the exact title
func (svr *ServerType) RemoveByTitle(resp http.ResponseWriter, req *http.Request) error {
	var todo types.TodoData
	err := decodeBody(req, &todo)
	if err != nil {
		respondErr(resp, req, http.StatusBadRequest, "Failed to decode body: ", err)
		return err
	}

	err = svr.DAO.DeleteByTitle(todo.Title)
	if err != nil {
		respondErr(resp, req, http.StatusInternalServerError, " Database error: ", err)
		return err
	}
	respond(resp, req, http.StatusOK, &types.ListStatus{
		Status: "Success",
		Info:   fmt.Sprintf("Todos with '%s' title removed", todo.Title),
	})
	return nil
}

//RemoveByPriority removes all todo items with the exact title
func (svr *ServerType) RemoveByPriority(resp http.ResponseWriter, req *http.Request) error {
	var todo types.TodoData
	err := decodeBody(req, &todo)
	if err != nil {
		respondErr(resp, req, http.StatusBadRequest, "Failed to decode body: ", err)
		return err
	}

	err = svr.DAO.DeleteByPriority(todo.Priority)
	if err != nil {
		respondErr(resp, req, http.StatusInternalServerError, " Database error: ", err)
		return err
	}
	respond(resp, req, http.StatusOK, &types.ListStatus{
		Status: "Success",
		Info:   fmt.Sprintf("Todos with '%d' priority removed", todo.Priority),
	})
	return nil
}

//RemoveByID removes the todo item with the given db id
func (svr *ServerType) RemoveByID(resp http.ResponseWriter, req *http.Request) error {
	var todo types.TodoData
	err := decodeBody(req, &todo)
	if err != nil {
		respondErr(resp, req, http.StatusBadRequest, "Failed to decode body: ", err)
		return err
	}

	err = svr.DAO.DeleteByID(todo.ID)
	if err != nil {
		respondErr(resp, req, http.StatusInternalServerError, " Database error: ", err)
		return err
	}
	respond(resp, req, http.StatusOK, &types.ListStatus{
		Status: "Success",
		Info:   fmt.Sprintf("Todo with '%d' id removed", todo.ID),
	})
	return nil
}

//ChangeTitle changes the title of the todo item by id
func (svr *ServerType) ChangeTitle(id int, resp http.ResponseWriter, req *http.Request) error {
	var todo types.TodoData
	err := decodeBody(req, &todo)
	if err != nil {
		respondErr(resp, req, http.StatusBadRequest, "Failed to decode body: ", err)
	}

	err = svr.DAO.UpdateTitle(todo.ID, todo.Title)
	if err != nil {
		respondErr(resp, req, http.StatusInternalServerError, " Database error: ", err)
		return err
	}
	respond(resp, req, http.StatusOK, &types.ListStatus{
		Status: "Success",
		Info:   fmt.Sprintf("Title changed to '%s' for id %d", todo.Title, id),
	})
	return nil
}

//ChangePriority changes the priority level of the todo item by id
func (svr *ServerType) ChangePriority(id int, resp http.ResponseWriter, req *http.Request) error {
	var todo types.TodoData
	err := decodeBody(req, &todo)
	if err != nil {
		respondErr(resp, req, http.StatusBadRequest, "Failed to decode body: ", err)
	}

	err = svr.DAO.UpdatePriority(todo.ID, todo.Priority)
	if err != nil {
		respondErr(resp, req, http.StatusInternalServerError, " Database error: ", err)
		return err
	}
	respond(resp, req, http.StatusOK, &types.ListStatus{
		Status: "Success",
		Info:   fmt.Sprintf("Priority changed to %d for id %d", todo.Priority, id),
	})
	return nil
}

//ChangeActive changes whether a given todo item is active or not
func (svr *ServerType) ChangeActive(id int, resp http.ResponseWriter, req *http.Request) error {
	var todo types.TodoData
	err := decodeBody(req, &todo)
	if err != nil {
		respondErr(resp, req, http.StatusBadRequest, "Failed to decode body: ", err)
	}

	err = svr.DAO.UpdateActive(todo.ID, todo.Active)
	if err != nil {
		respondErr(resp, req, http.StatusInternalServerError, " Database error: ", err)
		return err
	}
	respond(resp, req, http.StatusOK, &types.ListStatus{
		Status: "Success",
		Info:   fmt.Sprintf("Active changed to '%t' for id %d", todo.Active, id),
	})
	return nil
}
