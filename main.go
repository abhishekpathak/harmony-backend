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

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	log.Fatal(http.ListenAndServe(":" + port, loggedRouter))
}

func CheckError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func GetDbHandle() *sql.DB {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	err = db.Ping()
  	if err != nil {
    		panic(err)
  	}

	return db
}
