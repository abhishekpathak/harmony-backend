package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		handler := route.HandlerFunc
		handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

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
