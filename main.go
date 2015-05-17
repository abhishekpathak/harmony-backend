package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"time"
)

var SongID = 0

func play() {
	Truncate()
	Seed()
	ticker := time.NewTicker(time.Second)
	for _ = range ticker.C {
		go autoAdd()
		Refresh()
	}
}

func main() {
	go play()
	router := NewRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}

func GetDbHandle() *sql.DB {
	db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/songster")
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
