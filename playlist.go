package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Song struct {
	Id      int    `json:"id"`
	Videoid string `json:"videoid"`
	Name    string `json:"name"`
	Length  int    `json:"length"`
	Seek    int    `json:"seek"`
}

type Playlist []Song

func (s *Song) init(id int, videoid string, name string, length int, seek int) Song {
	return Song{
		Id:      id,
		Videoid: videoid,
		Name:    name,
		Length:  length,
		Seek:    seek,
	}
}

func createSong(videoid string, name string) Song {
	return Song{
		Id:      -1,
		Videoid: videoid,
		Name:    name,
		Length:  getDuration(videoid),
		Seek:    -5,
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
		s = s.init(id, videoid, name, length, seek)
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
	s = s.init(id, videoid, name, length, seek)
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
	db := GetDbHandle()
	defer db.Close()
	err := db.QueryRow("SELECT * FROM playlist ORDER BY id DESC LIMIT 1").Scan(&id, &videoid, &name, &length, &seek)
	CheckError(err)
	var s = Song{}
	lastSong := s.init(id, videoid, name, length, seek)
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

func Add(query string) {
	searchResults := cleanup(Search(query))
	enqueue(searchResults[0], "client")
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
