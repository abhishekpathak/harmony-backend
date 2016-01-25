package main

import (
	"net/http"

	"github.com/unbxd/goji/web/middleware"
	"github.com/zenazn/goji/web"
)

func RadioRoutes(route *web.Mux) *web.Mux {
	var radio *radioH
	route.Use(middleware.SubRouter)
	route.Use(setCORSmw)
	route.Use(setJsonContentType)
	route.Get("playlist/currentlyplaying", radio.currentlyPlaying)
	route.Get("playlist", radio.playlist)
	route.Post("add", radio.add)
	route.Post("skip", radio.skip)
	route.Get("query", radio.query)
	route.Get("lastsong", radio.lastSong)
	return route
}

func LibraryRoutes(route *web.Mux) *web.Mux {
	var lib *libH
	route.Use(middleware.SubRouter)
	route.Use(setCORSmw)
	route.Use(setJsonContentType)
	route.Get("update", lib.library)
	route.Get("songexists", lib.songExists)
	route.Get("get", lib.songFromLib)
	route.Get("updatelastplayed", lib.updateLibTimestamp)
	return route
}

func setCORSmw(context *web.C, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		handler.ServeHTTP(w, r)
	})
}

func setJsonContentType(context *web.C, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		handler.ServeHTTP(w, r)
	})
}
