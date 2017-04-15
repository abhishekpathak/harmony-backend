package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	for _, route := range routes {
		router.HandleFunc(route.Path, route.HandlerFunc).Methods(route.Method)
	}
	log.Fatal(http.ListenAndServe(":25404", router))
}

func CheckError(err error) {
	if err != nil {
		panic(err.Error())
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
