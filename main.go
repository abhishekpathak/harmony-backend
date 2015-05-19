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
