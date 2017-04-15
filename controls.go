package main

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/HarmonyProject/songster/musicservice"
)

func getVideoid(youtubeLink string) string {
	u, err := url.Parse(youtubeLink)
	CheckError(err)
	videoid := u.Query().Get("v")
	return videoid
}

func getQueryResults(query string) []musicservice.Song {
	var songs []musicservice.Song
	if strings.Contains(query, "www.youtube.com/watch?v=") {
		song := musicservice.CreateSong(getVideoid(query))
		if song.Length != -1 {
			songs = append(songs, song)
		}
	} else {
		songs = musicservice.Search(query)
	}
	return songs
}

func UpdateLibrary(form url.Values) bool {
	status := false
	var song musicservice.LibSong
	song.Videoid = form.Get("songvideoid")
	song.Artist = form.Get("songartist")
	song.Track = form.Get("songtrack")
	song.Rating, _ = strconv.Atoi(form.Get("songrating"))
	if form.Get("songfav") == "0" {
		song.Fav = false
	} else {
		song.Fav = true
	}

	var user musicservice.User
	user.Name = form.Get("username")
	user.Id = form.Get("userid")

	operation := form.Get("operation")

	if operation == "add" {
		status = addToLibrary(song, user)
	} else if operation == "remove" {
		status = removeFromLibrary(song, user)
	}

	return status
}
