package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/shale/go/client"
	"github.com/shale/go/data"
	"github.com/shale/go/service"
	log "github.com/sirupsen/logrus"
)

func init() {
	//Set up logging
	log.SetFormatter(&log.TextFormatter{DisableColors: true, FullTimestamp: true})
	log.SetLevel(log.DebugLevel)
}

func main() {
	log.Info("Starting main service")

	//Set env vars
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}
	test := false
	if os.Getenv("TEST") == "true" {
		test = true
	}

	//Set up db connection
	db, err := sql.Open("mysql", "root:root@tcp(db:3306)/sys?parseTime=true")
	//db, err := sql.Open("mysql", "root@tcp(localhost:3306)/sys?parseTime=true")	//Use for local testing outside of docker
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	dao := &data.StoreType{DAO: db}
	svc := &service.ServerType{DAO: dao}

	//Run a simple test client
	go func() {
		time.Sleep(time.Second * 10)
		if test {
			fmt.Printf("Running client\n")
			client.Run()
		}
	}()

	//Instantiate server and multiplexer, register endpoints, and start listening
	mux := http.NewServeMux()
	mux.HandleFunc("/todo/", svc.HandleTodos)
	log.Infof("Starting API on port %s", port)
	log.Fatal(http.ListenAndServe(port, mux))

	log.Info("Ending service")
}
