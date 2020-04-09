package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/bdlm/log"
	"github.com/shale/go/data"
	"github.com/shale/go/types"
)

//Server is the interface that defines the CRUD API object
type Server interface {
	HandleTodos(resp http.ResponseWriter, req *http.Request)
	AddTodo(name string) error
	GetTodos(name string, resp http.ResponseWriter, req *http.Request) error
	GetActives(active bool, resp http.ResponseWriter, req *http.Request) error
	GetTodosByPriority(priority int, resp http.ResponseWriter, req *http.Request) error
	GetTodosByCategory(category string, resp http.ResponseWriter, req *http.Request) error
	GetTodosByID(id int, resp http.ResponseWriter, req *http.Request) error
	RemoveByTitle(resp http.ResponseWriter, req *http.Request) error
	RemoveByPriority(resp http.ResponseWriter, req *http.Request) error
	RemoveInactive() error
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
	pathArgs := strings.Split(strings.Trim(req.URL.Path, "/"), "/")
	numArgs := len(pathArgs)
	fmt.Printf("\n%s: PATHARGS %+v\n", time.Now().String(), pathArgs)
	if numArgs < 2 {
		respondHTTPErr(resp, req, http.StatusBadRequest)
		return
	}
	name := pathArgs[1]
	switch req.Method {
	case "GET":
		var err error
		if numArgs == 2 {
			err = svr.GetTodos(name, resp, req)
			if err != nil {
				respondErr(resp, req, http.StatusBadRequest, " GET Error: ", err)
			}
			return
		} else if numArgs < 4 {
			respondHTTPErr(resp, req, http.StatusBadRequest)
			return
		}
		switch pathArgs[2] {
		case "active":
			active, err := strconv.ParseBool(pathArgs[3])
			if err == nil {
				err = svr.GetActives(active, name, resp, req)
			}
		case "highs":
			priority, err := strconv.Atoi(pathArgs[3])
			if err == nil {
				err = svr.GetTodosByPriority(priority, name, resp, req)
			}
		case "cat":
			err = svr.GetTodosByCategory(pathArgs[3], name, resp, req)
		case "id":
			id, err := strconv.Atoi(pathArgs[3])
			if err == nil {
				err = svr.GetTodosByID(id, name, resp, req)
			}
		}
		if err != nil {
			respondErr(resp, req, http.StatusBadRequest, " GET Error: ", err)
			log.Errorf("GET Error: %v", err)
		}
		return
	case "POST":
		var err error
		if numArgs == 3 && pathArgs[2] == "add" {
			err = svr.AddTodo(name, resp, req)
			if err != nil {
				respondErr(resp, req, http.StatusBadRequest, " POST Error: ", err)
				log.Errorf("POST Error: %v", err)
			}
			return
		} else if numArgs < 4 {
			respondHTTPErr(resp, req, http.StatusBadRequest)
			return
		}
		switch pathArgs[2] {
		case "ctitle":
			id, err := strconv.Atoi(pathArgs[3])
			if err == nil {
				err = svr.ChangeTitle(id, name, resp, req)
			}
		case "cpri":
			id, err := strconv.Atoi(pathArgs[3])
			if err == nil {
				err = svr.ChangePriority(id, name, resp, req)
			}
		case "cactive":
			id, err := strconv.Atoi(pathArgs[3])
			if err == nil {
				err = svr.ChangeActive(id, name, resp, req)
			}
		}
		if err != nil {
			respondErr(resp, req, http.StatusBadRequest, " POST Error: ", err)
			log.Errorf("POST Error: %v", err)
		}
		return
	case "DELETE":
		var err error
		if numArgs < 3 {
			respondHTTPErr(resp, req, http.StatusBadRequest)
			return
		}
		switch pathArgs[2] {
		case "rmtitle":
			err = svr.RemoveByTitle(name, resp, req)
		case "rmpri":
			err = svr.RemoveByPriority(name, resp, req)
		case "rminactive":
			err = svr.RemoveInactive(name, resp, req)
		case "rmid":
			err = svr.RemoveByID(name, resp, req)
		}
		if err != nil {
			respondErr(resp, req, http.StatusBadRequest, " DELETE Error: ", err)
			log.Errorf("DELETE Error: %v", err)
		}
		return
	default:
		respondHTTPErr(resp, req, http.StatusBadRequest)
		return
	}
}

//AddTodo adds a new to do item to the to do list
func (svr *ServerType) AddTodo(name string, resp http.ResponseWriter, req *http.Request) error {
	var todo types.TodoData
	err := decodeBody(req, &todo)
	if err != nil {
		return err
	}
	todo.Name = name
	err = svr.DAO.InsertTodo(todo)
	if err != nil {
		return err
	}
	respond(resp, req, http.StatusOK, &types.ListStatus{
		Status: "Success",
		Info:   fmt.Sprintf("Title '%s' added", todo.Title),
	})
	return nil
}

//GetTodos returns all of the todo items from the list
func (svr *ServerType) GetTodos(name string, resp http.ResponseWriter, req *http.Request) error {
	result, err := svr.DAO.SelectAllTodos(name)
	if err != nil {
		return err
	}
	if result == nil {
		respond(resp, req, http.StatusNoContent, &result)
	} else {
		respond(resp, req, http.StatusOK, &result)
	}
	return nil
}

//GetActives returns all of the todo items from the list
func (svr *ServerType) GetActives(active bool, name string, resp http.ResponseWriter, req *http.Request) error {
	result, err := svr.DAO.SelectActives(active, name)
	if err != nil {
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
func (svr *ServerType) GetTodosByPriority(priority int, name string, resp http.ResponseWriter, req *http.Request) error {
	var result []types.TodoData
	var err error
	if priority == 0 {
		result, err = svr.DAO.SelectNonPriority(name)
	} else {
		result, err = svr.DAO.SelectByPriority(priority, name)
	}
	if err != nil {
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
func (svr *ServerType) GetTodosByCategory(category string, name string, resp http.ResponseWriter, req *http.Request) error {
	result, err := svr.DAO.SelectByCategory(category, name)
	if err != nil {
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
func (svr *ServerType) GetTodosByID(id int, name string, resp http.ResponseWriter, req *http.Request) error {
	result, err := svr.DAO.SelectByID(id, name)
	if err != nil {
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
func (svr *ServerType) RemoveByTitle(name string, resp http.ResponseWriter, req *http.Request) error {
	var todo types.TodoData
	err := decodeBody(req, &todo)
	if err != nil {
		return err
	}

	err = svr.DAO.DeleteByTitle(todo.Title, name)
	if err != nil {
		return err
	}
	respond(resp, req, http.StatusOK, &types.ListStatus{
		Status: "Success",
		Info:   fmt.Sprintf("Todos with '%s' title removed", todo.Title),
	})
	return nil
}

//RemoveByPriority removes all todo items with the exact title
func (svr *ServerType) RemoveByPriority(name string, resp http.ResponseWriter, req *http.Request) error {
	var todo types.TodoData
	err := decodeBody(req, &todo)
	if err != nil {
		return err
	}

	err = svr.DAO.DeleteByPriority(todo.Priority, name)
	if err != nil {
		return err
	}
	respond(resp, req, http.StatusOK, &types.ListStatus{
		Status: "Success",
		Info:   fmt.Sprintf("Todos with '%d' priority removed", todo.Priority),
	})
	return nil
}

//RemoveInactive removes the todo item with the given db id
func (svr *ServerType) RemoveInactive(name string, resp http.ResponseWriter, req *http.Request) error {
	err := svr.DAO.DeleteInactive(name)
	if err != nil {
		return err
	}
	respond(resp, req, http.StatusOK, &types.ListStatus{
		Status: "Success",
		Info:   fmt.Sprintf("Todos with inactive status removed"),
	})
	return nil
}

//RemoveByID removes the todo item with the given db id
func (svr *ServerType) RemoveByID(name string, resp http.ResponseWriter, req *http.Request) error {
	var todo types.TodoData
	err := decodeBody(req, &todo)
	if err != nil {
		return err
	}

	err = svr.DAO.DeleteByID(todo.ID, name)
	if err != nil {
		return err
	}
	respond(resp, req, http.StatusOK, &types.ListStatus{
		Status: "Success",
		Info:   fmt.Sprintf("Todo with '%d' id removed", todo.ID),
	})
	return nil
}

//ChangeTitle changes the title of the todo item by id
func (svr *ServerType) ChangeTitle(id int, name string, resp http.ResponseWriter, req *http.Request) error {
	var todo types.TodoData
	err := decodeBody(req, &todo)
	if err != nil {
		return err
	}

	err = svr.DAO.UpdateTitle(id, todo.Title, name)
	if err != nil {
		return err
	}
	respond(resp, req, http.StatusOK, &types.ListStatus{
		Status: "Success",
		Info:   fmt.Sprintf("Title changed to '%s' for id %d", todo.Title, id),
	})
	return nil
}

//ChangePriority changes the priority level of the todo item by id
func (svr *ServerType) ChangePriority(id int, name string, resp http.ResponseWriter, req *http.Request) error {
	var todo types.TodoData
	err := decodeBody(req, &todo)
	if err != nil {
		return err
	}

	err = svr.DAO.UpdatePriority(id, todo.Priority, name)
	if err != nil {
		return err
	}
	respond(resp, req, http.StatusOK, &types.ListStatus{
		Status: "Success",
		Info:   fmt.Sprintf("Priority changed to %d for id %d", todo.Priority, id),
	})
	return nil
}

//ChangeActive changes whether a given todo item is active or not
func (svr *ServerType) ChangeActive(id int, name string, resp http.ResponseWriter, req *http.Request) error {
	var todo types.TodoData
	err := decodeBody(req, &todo)
	if err != nil {
		return err
	}

	err = svr.DAO.UpdateActive(id, todo.Active, name)
	if err != nil {
		return err
	}
	respond(resp, req, http.StatusOK, &types.ListStatus{
		Status: "Success",
		Info:   fmt.Sprintf("Active changed to '%t' for id %d", todo.Active, id),
	})
	return nil
}
