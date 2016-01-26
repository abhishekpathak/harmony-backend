package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/zenazn/goji/web"
)

func main() {
	go play()
	go autoAdd()

	router := web.New()

	router.Handle("/radio/*", RadioRoutes(web.New())) // Radio Sub-Router
	router.Handle("/library", http.RedirectHandler("/library/update", http.StatusMovedPermanently))
	router.Handle("/library/*", LibraryRoutes(web.New())) // Radio Sub-Router

	err := http.ListenAndServe(":25404", router)
	if err != nil {
		log.Fatal("Error Starting Server", err)
	}
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
