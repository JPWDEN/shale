package main

import (
	"database/sql"
	"net/http"
	"os"

	"github.com/shale/go/data"
	"github.com/shale/go/service"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	log.SetLevel(log.DebugLevel)
}

func main() {
	log.Info("Starting main service")

	//Logging?

	//Set env vars
	port := os.Getenv("PORT")
	addr := os.Getenv("ADDRESS")
	if port == "" || addr == "" {
		port = ":8080"
		addr = "localhost"
	}
	test := false
	if os.Getenv("TEST") == "true" {
		test = true
	}

	//Set up db connection
	db, err := sql.Open("mysql", "root@tcp(db:3306)/sys?parseTime=true")
	//db, err := sql.Open("mysql", "root@tcp(localhost:3306)/sys?parseTime=true")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	dao := &data.StoreType{DAO: db}
	svc := &service.ServerType{DAO: dao}

	//Instantiate server and multiplexer, register endpoints, and start listening
	mux := http.NewServeMux()
	mux.HandleFunc("/todo/", svc.HandleTodos)
	log.Infof("Starting API on %s", addr)
	log.Fatal(http.ListenAndServe(port, mux))

	//Test client
	go func() {
		if !test {
			return
		}
		//for {
		//}
	}()

	log.Info("Ending service")
}
