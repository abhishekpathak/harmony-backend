package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"time"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		handler := route.HandlerFunc
		handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			f, err := os.OpenFile("/Users/abhishek.p/logs/songster/root.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				fmt.Printf("error opening file: %v", err)
			}
			defer f.Close()
			log.SetOutput(f)

			start := time.Now()
			handler.ServeHTTP(w, r)
			log.Printf(
				"%s\t%s\t%s\t%s",
				r.Method,
				r.RequestURI,
				route,
				time.Since(start),
			)
		})

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)

	}
	return router
}
