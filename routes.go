package main

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"CurrentlyPlaying",
		"GET",
		"/playlist/currentlyplaying",
		CurrentlyPlayingHandler,
	},
	Route{
		"Playlist",
		"GET",
		"/playlist",
		PlaylistHandler,
	},
	Route{
		"Add",
		"POST",
		"/add",
		AddHandler,
	},
	Route{
		"Skip",
		"POST",
		"/skip",
		SkipHandler,
	},
	Route{
		"QueryResults",
		"GET",
		"/query",
		QueryHandler,
	},
	Route{
		"LastSong",
		"GET",
		"/lastsong",
		LastSongHandler,
	},
}
