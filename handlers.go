package main

import (
	"encoding/json"
	"net/http"
)

func CurrentlyPlayingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	song := CurrentlyPlaying()
	json.NewEncoder(w).Encode(song)
}

func PlaylistHandler(w http.ResponseWriter, r *http.Request) {
	playlist := GetPlaylist()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(playlist)
}

func AddHandler(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("q")
	user := r.FormValue("user")
	status := UserAdd(query, user)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if status == true {
		w.Write([]byte("{\"status\":\"success\"}"))
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "{\"status\":\"error\"}", http.StatusBadRequest)
	}
}

func SkipHandler(w http.ResponseWriter, r *http.Request) {
	Skip()
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write([]byte("{\"status\":\"success\"}"))
}

func QueryHandler(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("q")
	matchedResults := getQueryResults(query)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(matchedResults)
}

func LastSongHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.FormValue("userid")
	lastSong := getLastPlaying(userId)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(lastSong)
}
