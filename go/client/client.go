package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/bdlm/log"
)

//Get todos:  curl -vv 73.78.155.49:8080/todo/tom
func runGetTodos() {
	route := "http://localhost:8080/todo/jpw"
	req, err := http.NewRequest("GET", route, nil)
	if err != nil {
		log.Infof("Error %v", err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Infof("Error %v", err)
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	log.Infof(string(body))
}

//Get Active:  curl -vv 73.78.155.49:8080/todo/tom/active/1
func runGetActives() {
	route := "http://localhost:8080/todo/jpw/active/1"
	req, err := http.NewRequest("GET", route, nil)
	if err != nil {
		log.Infof("Error %v", err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Infof("Error %v", err)
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	log.Infof(string(body))
}

//ByPriority:  curl -vv 73.78.155.49:8080/todo/tom/highs/6
func runGetByPriority() {
	route := "http://localhost:8080/todo/jpw/pri/2"
	req, err := http.NewRequest("GET", route, nil)
	if err != nil {
		log.Infof("Error %v", err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Infof("Error %v", err)
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	log.Infof(string(body))
}

//ByCategory:  curl -vv 73.78.155.49:8080/todo/tom/cat/shopping
func runGetByCategory() {
	route := "http://localhost:8080/todo/jpw/cat/shopping"
	req, err := http.NewRequest("GET", route, nil)
	if err != nil {
		log.Infof("Error %v", err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Infof("Error %v", err)
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	log.Infof(string(body))
}

//ByID:  curl -vv 73.78.155.49:8080/todo/tom/id/4
func runGetByID() {
	route := "http://localhost:8080/todo/jpw/id/4"
	req, err := http.NewRequest("GET", route, nil)
	if err != nil {
		log.Infof("Error %v", err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Infof("Error %v", err)
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	log.Infof(string(body))
}

//Add:
//curl -v -X POST localhost:8080/todo/tom/add --data "{\"title\": \"#2\", \"item_priority\": 5, \"category\": \"shopping\"}"
//curl -v -X POST localhost:8080/todo/tom/add --data "{\"title\": \"Brand New\", \"body\": \"This is fun\", \"item_priority\": 8, \"category\": \"entertainment\"}"
func runAdd() {
	route := "http://localhost:8080/todo/jpw/add"
	payload := map[string]interface{}{
		"title":         "Test Add",
		"item_priority": 3,
		"category":      "Tests",
	}
	byteMap, err := json.Marshal(payload)
	if err != nil {
		log.Errorf("Error: %v", err)
	}
	resp, err := http.Post(route, "application/json", bytes.NewBuffer(byteMap))
	if err != nil {
		log.Errorf("Error: %v", err)
	}
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	log.Infof("result: %v", result)
}

//Change Title:  curl -v -X POST localhost:8080/todo/tom/ctitle/4 --data "{\"title\": \"Changed TITLE\"}"
func runChangeTitle() {
	route := "http://localhost:8080/todo/jpw/ctitle/2"
	payload := map[string]interface{}{
		"title": "Changed Title",
	}
	byteMap, err := json.Marshal(payload)
	if err != nil {
		log.Errorf("Error: %v", err)
	}
	resp, err := http.Post(route, "application/json", bytes.NewBuffer(byteMap))
	if err != nil {
		log.Errorf("Error: %v", err)
	}
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	log.Infof("result: %v", result)
}

//Change Priority:  curl -v -X POST localhost:8080/todo/tom/cpri/4 --data "{\"item_priority\": 2}"
func runChangePriority() {
	route := "http://localhost:8080/todo/jpw/cpri/2"
	payload := map[string]interface{}{
		"item_priority": 9,
	}
	byteMap, err := json.Marshal(payload)
	if err != nil {
		log.Errorf("Error: %v", err)
	}
	resp, err := http.Post(route, "application/json", bytes.NewBuffer(byteMap))
	if err != nil {
		log.Errorf("Error: %v", err)
	}
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	log.Infof("result: %v", result)
}

//Change Active:  curl -v -X POST localhost:8080/todo/tom/cactive/2 --data "{\"active\": false}"
func runChangeActive() {
	route := "http://localhost:8080/todo/jpw/cactive/2"
	payload := map[string]interface{}{
		"active": false,
	}
	byteMap, err := json.Marshal(payload)
	if err != nil {
		log.Errorf("Error: %v", err)
	}
	resp, err := http.Post(route, "application/json", bytes.NewBuffer(byteMap))
	if err != nil {
		log.Errorf("Error: %v", err)
	}
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	log.Infof("result: %v", result)
}

//Remove Title:  curl -v -X DELETE localhost:8080/todo/tom/rmtitle --data "{\"title\": \"#2\"}"
func runRemoveTitle() {
	route := "http://localhost:8080/todo/jpw/rmtitle"
	payload := map[string]interface{}{
		"title": "#2",
	}
	byteMap, err := json.Marshal(payload)
	if err != nil {
		log.Errorf("Error: %v", err)
	}
	req, err := http.NewRequest("DELETE", route, bytes.NewBuffer(byteMap))
	if err != nil {
		log.Infof("Error %v", err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Infof("Error %v", err)
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	log.Infof(string(body))
}

//Remove Priority:  curl -v -X DELETE localhost:8080/todo/tom/rmpri --data "{\"item_priority\": 0}"
func runRemovePriority() {
	route := "http://localhost:8080/todo/jpw/rmpri"
	payload := map[string]interface{}{
		"item_priority": 9,
	}
	byteMap, err := json.Marshal(payload)
	if err != nil {
		log.Errorf("Error: %v", err)
	}
	req, err := http.NewRequest("DELETE", route, bytes.NewBuffer(byteMap))
	if err != nil {
		log.Infof("Error %v", err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Infof("Error %v", err)
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	log.Infof(string(body))
}

//Remove id:  curl -v -X DELETE localhost:8080/todo/tom/rmid --data "{\"id\": 8}"
func runRemoveID() {
	route := "http://localhost:8080/todo/jpw/rmid"
	payload := map[string]interface{}{
		"id": 10,
	}
	byteMap, err := json.Marshal(payload)
	if err != nil {
		log.Errorf("Error: %v", err)
	}
	req, err := http.NewRequest("DELETE", route, bytes.NewBuffer(byteMap))
	if err != nil {
		log.Infof("Error %v", err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Infof("Error %v", err)
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	log.Infof(string(body))
}

// Run runs each of the functional tests in this package
func Run() {
	runGetTodos()
	//runGetActives()
	//runGetByPriority()
	//runGetByCategory()
	//runGetByID()
	//runAdd()
	//runChangeTitle()
	//runChangePriority()
	//runChangeActive()
	//runRemoveTitle()
	//runRemovePriority()
	//runRemoveID()
	runGetTodos()
}
