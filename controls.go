package main

import (
	"fmt"
	"github.com/HarmonyProject/songster/musicservice"
	"net/url"
	"strings"
	"time"
)

func play() {
	if playlistSize() == 0 {
		seed()
	}
	ticker := time.NewTicker(time.Second)
	for _ = range ticker.C {
		refresh()
	}
}

func enqueue(s musicservice.Song, agent string) {
	addToPlaylist(s, agent)
	UpdateSongdetails(s)
}

func CurrentlyPlaying() musicservice.Song {
	s := getSong(firstSongId())
	return s
}

func refresh() {
	s := CurrentlyPlaying()
	if s.Seek < s.Length {
		updateSeek(s.Id)
		fmt.Printf("\r%d/%d - %s  ", s.Seek, s.Length, s.Name)
	} else {
		remove(s)
		refresh()
	}
}

func Skip() {
	s := CurrentlyPlaying()
	remove(s)
	refresh()
}

func seed() {
	seedQuery := "tum se hi"
	searchResults := musicservice.Search(seedQuery)
	seedSong := searchResults[0]
	clearPlaylist()
	enqueue(seedSong, "system")
}

func GetPlaylist() []musicservice.Song {
	var playlist []musicservice.Song
	ids := currentPlaylistIds()
	for _, id := range ids {
		playlist = append(playlist, getSong(id))
	}
	return playlist
}

func getLastSong() musicservice.Song {
	s := getSong(lastSongId())
	return s
}

func remove(s musicservice.Song) {
	removeFromPlaylist(s.Id)
}

func getVideoid(youtubeLink string) string {
	u, err := url.Parse(youtubeLink)
	if err != nil {
		fmt.Println("unable to parse URL")
	}
	videoid := u.Query().Get("v")
	return videoid
}

func UserAdd(query string, user string) bool {
	if strings.Contains(query, "www.youtube.com/watch?v=") {
		song := musicservice.CreateSong(getVideoid(query))
		if song.Length == -1 || user == "" {
			return false
		} else {
			enqueue(song, user)
			return true
		}
	}
	searchResults := musicservice.Search(query)
	if len(searchResults) == 0 {
		return false
	}
	for i := range searchResults {
		if searchResults[i].Name == query {
			enqueue(searchResults[i], user)
			return true
		}
	}
	enqueue(searchResults[0], user)
	return true
}

func autoAdd() {
	ticker := time.NewTicker(time.Second * 5)
	for _ = range ticker.C {
		c := CurrentlyPlaying()
		timeRemaining := c.Length - c.Seek
		if playlistSize() == 1 && timeRemaining < 30 {
			newSong := musicservice.Recommend(getLastSong())
			enqueue(newSong, "system")
		}
	}
}
