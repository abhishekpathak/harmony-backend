package main

import (
	"encoding/json"
	"github.com/HarmonyProject/songster/musicservice"
	"net/http"
	"strconv"
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

func LibraryHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	status := UpdateLibrary(r.Form)
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

func SongExistsHandler(w http.ResponseWriter, r *http.Request) {
	userid := r.FormValue("userid")
	videoid := r.FormValue("videoid")
	status := songExistsInLibrary(userid, videoid)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"status\":" + strconv.FormatBool(status) + "}"))
}

func SongFromLibraryHandler(w http.ResponseWriter, r *http.Request) {
	userid := r.FormValue("userid")
	fav := r.FormValue("fav")
	var libSongObj musicservice.LibSong
	if fav == "true" {
		libSongObj = favSongFromLibrary(userid)
	} else {
		libSongObj = randomSongFromLibrary(userid)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(libSongObj)
}

func UpdateLibraryTimestampHandler(w http.ResponseWriter, r *http.Request) {
	userid := r.FormValue("userid")
	videoid := r.FormValue("videoid")
	status := updateLastPlayedTimestamp(userid, videoid)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if status == true {
		w.Write([]byte("{\"status\":\"success\"}"))
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "{\"status\":\"error\"}", http.StatusBadRequest)
	}
}
