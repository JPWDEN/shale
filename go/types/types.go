package types

import "github.com/go-sql-driver/mysql"

//TodoData is the JSON-relatable object used for API call
type TodoData struct {
	Name        string         `json:"acct_name"`
	Title       string         `json:"title"`
	Body        string         `json:"body"`
	Category    string         `json:"category"`
	Priority    int            `json:"item_priority"`
	PublishDate mysql.NullTime `json:"publish_date"`
	Active      bool           `json:"active"`
	ID          int            `json:"id"`
}

//ListStatus prides a status response for changes made to the todo list
type ListStatus struct {
	Status string `json:"status"`
	Info   string `json:"info"`
}
