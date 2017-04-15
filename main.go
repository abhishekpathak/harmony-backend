package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
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

func GetDbHandle() *sql.DB {
	DB_HOST := "localhost"
	DB_PORT := "3306"
	DB_NAME := "songster"
	DSN := "songster:songster@tcp(" + DB_HOST + ":" + DB_PORT + ")/" + DB_NAME
	db, err := sql.Open("mysql", DSN)
	if err != nil {
		panic(err.Error())
	}
	return db
}
