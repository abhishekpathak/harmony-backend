package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	for _, route := range routes {
		router.HandleFunc(route.Path, route.HandlerFunc).Methods(route.Method)
	}
	loggedRouter := handlers.LoggingHandler(os.Stdout, router)
	log.Fatal(http.ListenAndServe(":25404", loggedRouter))
}

func CheckError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

const (
  host     = "localhost"
  port     = 5432
  user     = "harmony"
  password = "harmony"
  dbname   = "harmony"
)

func GetDbHandle() *sql.DB {
	DSN := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", DSN)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
  	if err != nil {
    		panic(err)
  	}

	return db
}
