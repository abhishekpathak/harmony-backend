package main

import "net/http"

type Route struct {
	Name        string
	Method      string
	Path        string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"QueryResults",
		"GET",
		"/query",
		QueryHandler,
	},
	Route{
		"UpdateLibrary",
		"GET",
		"/library",
		LibraryHandler,
	},
	Route{
		"Song Exists",
		"GET",
		"/library/songexists",
		SongExistsHandler,
	},
	Route{
		"PlayFromLibrary",
		"GET",
		"/library/get",
		SongFromLibraryHandler,
	},
	Route{
		"UpdateLibraryTimestamp",
		"GET",
		"/library/updatelastplayed",
		UpdateLibraryTimestampHandler,
	},
	Route{
		"GetRecommendedSongs",
		"GET",
		"/recommendations",
		RecommendationHandler,
	},
}
