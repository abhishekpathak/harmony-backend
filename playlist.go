package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type Song struct {
	Id      int    `json:"id"`
	Videoid string `json:"videoid"`
	Name    string `json:"name"`
	Length  int    `json:"length"`
	Seek    int    `json:"seek"`
	AddedBy string `json:"added_by"`
}

type Playlist []Song

func (s *Song) init(id int, videoid string, name string, length int, seek int, addedBy string) Song {
	return Song{
		Id:      id,
		Videoid: videoid,
		Name:    name,
		Length:  length,
		Seek:    seek,
		AddedBy: addedBy,
	}
}

func createSong(videoid string, name string) Song {
	return Song{
		Id:      -1,
		Videoid: videoid,
		Name:    name,
		Length:  getDuration(videoid),
		Seek:    -5,
		AddedBy: "system",
	}
}

func Truncate() {
	db := GetDbHandle()
	defer db.Close()
	stmt, err := db.Prepare("TRUNCATE playlist")
	CheckError(err)
	_, err = stmt.Exec()
	CheckError(err)
}

func cleanup(results []Song) []Song {
	var cleanedResults []Song
	for i := range results {
		if results[i].Length != -1 {
			cleanedResults = append(cleanedResults, results[i])
		}
	}
	return cleanedResults
}

func Seed() {
	seedQuery := "tum se hi"
	searchResults := cleanup(Search(seedQuery))
	seedSong := searchResults[0]
	Truncate()
	enqueue(seedSong, "system")
}

func GetPlaylist() Playlist {
	var id int
	var videoid string
	var name string
	var length int
	var seek int
	var addedBy string
	playlist := []Song{}
	db := GetDbHandle()
	defer db.Close()
	rows, err := db.Query("SELECT  * from  playlist order by id")
	CheckError(err)
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &videoid, &name, &length, &seek, &addedBy)
		var s = Song{}
		s = s.init(id, videoid, name, length, seek, addedBy)
		CheckError(err)
		playlist = append(playlist, s)
	}
	return playlist
}

func CurrentlyPlaying() Song {
	var id int
	var videoid string
	var name string
	var length int
	var seek int
	var addedBy string
	db := GetDbHandle()
	defer db.Close()
	err := db.QueryRow("SELECT * FROM playlist ORDER BY id ASC LIMIT 1").Scan(&id, &videoid, &name, &length, &seek, &addedBy)
	CheckError(err)
	var s = Song{}
	s = s.init(id, videoid, name, length, seek, addedBy)
	return s
}

func updateSeek(s Song) {
	db := GetDbHandle()
	defer db.Close()
	stmt, err := db.Prepare("UPDATE playlist SET seek = seek + 1 WHERE id = ?")
	CheckError(err)
	_, err = stmt.Exec(s.Id)
	CheckError(err)
}

func getLastSong() Song {
	var id int
	var videoid string
	var name string
	var length int
	var seek int
	var addedBy string
	db := GetDbHandle()
	defer db.Close()
	err := db.QueryRow("SELECT * FROM playlist ORDER BY id DESC LIMIT 1").Scan(&id, &videoid, &name, &length, &seek, &addedBy)
	CheckError(err)
	var s = Song{}
	lastSong := s.init(id, videoid, name, length, seek, addedBy)
	return lastSong
}

func enqueue(s Song, agent string) {
	db := GetDbHandle()
	defer db.Close()
	stmt, err := db.Prepare("INSERT INTO playlist (videoid, name, length, seek, added_by) VALUES (?, ?, ?, ?, ?)")
	CheckError(err)
	_, err = stmt.Exec(s.Videoid, s.Name, s.Length, s.Seek, agent)
	CheckError(err)
}

func remove(s Song) {
	db := GetDbHandle()
	defer db.Close()
	stmt, err := db.Prepare("DELETE FROM playlist WHERE id = ?")
	CheckError(err)
	_, err = stmt.Exec(s.Id)
	CheckError(err)
}

func Size() int {
	var size int
	db := GetDbHandle()
	defer db.Close()
	err := db.QueryRow("SELECT count(*) FROM playlist").Scan(&size)
	CheckError(err)
	return size
}

func Add(query string, user string) bool {
	if strings.Contains(query, "www.youtube.com/watch?v=") {
		query = strings.Replace(query, "https://www.youtube.com/watch?v=", "", 1)
		query = strings.Replace(query, "http://www.youtube.com/watch?v=", "", 1)
		query = strings.Replace(query, "www.youtube.com/watch?v=", "", 1)
		song := createSong(query, GetName(query))
		if song.Length == -1 || user == "" {
			return false
		} else {
			enqueue(song, user)
			return true
		}
	}
	searchResults := cleanup(Search(query))
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
		if Size() == 1 && timeRemaining < 30 {
			newSong := recommend(getLastSong())
			enqueue(newSong, "system")
		}
	}
}

func recommend(s Song) Song {
	var recommendedSong Song
	recommendations := cleanup(Recommend(s.Videoid))
	if len(recommendations) < 6 {
		seedQuery := "tum se hi"
		searchResults := cleanup(Search(seedQuery))
		recommendedSong = searchResults[0]
	} else {
		songindex := rand.Intn(len(recommendations)-3) + 3
		recommendedSong = recommendations[songindex]
	}
	return recommendedSong
}

func Refresh() {
	s := CurrentlyPlaying()
	if s.Seek < s.Length {
		updateSeek(s)
		fmt.Println(s.Videoid, "   ", s.Seek, "          ", GetPlaylist())
	} else {
		remove(s)
		go PostToSlack("#nowplaying " + CurrentlyPlaying().Name)
		Refresh()
	}
}

func Skip() {
	s := CurrentlyPlaying()
	PostToSlack("#skipped " + CurrentlyPlaying().Name + ". Don't do this :rage:")
	remove(s)
	go PostToSlack("#nowplaying " + CurrentlyPlaying().Name)
	Refresh()
}
