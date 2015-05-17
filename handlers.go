package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func CurrentlyPlayingHandler(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	//w.WriteHeader(http.StatusOK)
	song := CurrentlyPlaying()
	fmt.Println(song)
	json.NewEncoder(w).Encode(song)
}

func PlaylistHandler(w http.ResponseWriter, r *http.Request) {
	playlist := GetPlaylist()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(playlist)
}
