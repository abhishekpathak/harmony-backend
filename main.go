package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

var SongID = 0

func play() {
	//Seed()
	go PostToSlack("#nowplaying " + CurrentlyPlaying().Name)
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

func PostToSlack(text string) {
	// constructing the URL
	apiUrl := "https://hooks.slack.com"
	resource := "/services/T04TS7W4P/B050S572D/T8pN7F6QSlI7I6PB90clC25I"
	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resource
	urlStr := fmt.Sprintf("%v", u)

	textPayload := `{"text": ` + text + `"}`
	data := url.Values{}
	data.Set("payload", textPayload)
	client := &http.Client{}
	fmt.Println(urlStr, textPayload)
	resp, err := client.PostForm(urlStr, data)
	CheckError(err)
	fmt.Println("Slack webhook responded with code : ", resp)
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
