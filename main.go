package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"os"
	"time"
)

var SongID = 0

func play() {
	Seed()
	ticker := time.NewTicker(time.Second)
	for _ = range ticker.C {
		Refresh()
	}
}

func main() {
	go play()
	go autoAdd()
	router := NewRouter()
	log.Fatal(http.ListenAndServe(":25404", router))
}

func GetDbHandle() *sql.DB {
	DB_HOST := os.Getenv("OPENSHIFT_MYSQL_DB_HOST")
	DB_PORT := os.Getenv("OPENSHIFT_MYSQL_DB_PORT")
	DB_NAME := "songster"
	DSN := "root@tcp(" + DB_HOST + ":" + DB_PORT + ")/" + DB_NAME
	db, err := sql.Open("mysql", DSN)
	CheckError(err)
	return db
}

func CheckError(err error) {
	if err != nil {
		handleError(err)
	}
}

func handleError(err error) {
	panic(err.Error())
}
