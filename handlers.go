package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/HarmonyProject/songster/musicservice"
)

// Radio Handlers
type radioH struct{}

func (*radioH) currentlyPlaying(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	song := CurrentlyPlaying()
	json.NewEncoder(w).Encode(song)
}

func (*radioH) playlist(w http.ResponseWriter, r *http.Request) {
	playlist := GetPlaylist()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(playlist)
}

func (*radioH) add(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("q")
	user := r.FormValue("user")
	status := UserAdd(query, user)
	if status == true {
		w.Write([]byte("{\"status\":\"success\"}"))
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "{\"status\":\"error\"}", http.StatusBadRequest)
	}
}

func (*radioH) skip(w http.ResponseWriter, r *http.Request) {
	Skip()
	w.Write([]byte("{\"status\":\"success\"}"))
}

func (*radioH) query(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("q")
	matchedResults := getQueryResults(query)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(matchedResults)
}

func (*radioH) lastSong(w http.ResponseWriter, r *http.Request) {
	userId := r.FormValue("userid")
	lastSong := getLastPlaying(userId)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(lastSong)
}

// Library Handlers

type libH struct{}

func (*libH) library(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	status := UpdateLibrary(r.Form)
	if status == true {
		w.Write([]byte("{\"status\":\"success\"}"))
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "{\"status\":\"error\"}", http.StatusBadRequest)
	}
}

func (*libH) songExists(w http.ResponseWriter, r *http.Request) {
	userid := r.FormValue("userid")
	videoid := r.FormValue("videoid")
	status := songExistsInLibrary(userid, videoid)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"status\":" + strconv.FormatBool(status) + "}"))
}

func (*libH) songFromLib(w http.ResponseWriter, r *http.Request) {
	userid := r.FormValue("userid")
	fav := r.FormValue("fav")
	var libSongObj musicservice.LibSong
	if fav == "true" {
		libSongObj = favSongFromLibrary(userid)
	} else {
		libSongObj = randomSongFromLibrary(userid)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(libSongObj)
}

func (*libH) updateLibTimestamp(w http.ResponseWriter, r *http.Request) {
	userid := r.FormValue("userid")
	videoid := r.FormValue("videoid")
	status := updateLastPlayedTimestamp(userid, videoid)
	if status == true {
		w.Write([]byte("{\"status\":\"success\"}"))
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "{\"status\":\"error\"}", http.StatusBadRequest)
	}
}
