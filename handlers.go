package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/HarmonyProject/songster/musicservice"
)

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

func QueryHandler(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("q")
	matchedResults := getQueryResults(query)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if len(matchedResults) == 0 {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(matchedResults)
	}
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
	var libSongObj musicservice.LibSong
	libSongObj = randomSongFromLibrary(userid)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if libSongObj.Videoid == "" {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(libSongObj)
	}
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

func RecommendationHandler(w http.ResponseWriter, r *http.Request) {
	videoid := r.FormValue("q")
	recommendedSongs := musicservice.RecommendMulti(videoid)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(recommendedSongs)
}
