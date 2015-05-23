package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func CurrentlyPlayingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	song := CurrentlyPlaying()
	fmt.Println(song)
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
	status := Add(query, user)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if status == true {
		http.Error(w, "{\"status\":\"error\"}", http.StatusOK)
	} else {
		http.Error(w, "{\"status\":\"error\"}", http.StatusBadRequest)
	}
}

func SkipHandler(w http.ResponseWriter, r *http.Request) {
	Skip()
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write([]byte("{\"status\":\"success\"}"))
}
