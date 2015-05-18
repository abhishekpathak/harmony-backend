package main

import (
	"fmt"
	"time"
)

type Song struct {
	Id      int    `json:"id"`
	Videoid string `json:"videoid"`
	Length  int    `json:"length"`
	Seek    int    `json:"seek"`
}

func createSong(videoid string) Song {
	return Song{
		Id:      -1,
		Videoid: videoid,
		Length:  getDuration(videoid),
		Seek:    0,
	}
}

type Playlist []Song

func Truncate() {
	db := GetDbHandle()
	defer db.Close()
	stmt, err := db.Prepare("TRUNCATE playlist")
	CheckError(err)
	_, err = stmt.Exec()
	CheckError(err)
}

func Seed() {
	seedQuery := "David Bowie Heroes"
	searchResults := Search(seedQuery)
	seedSong := searchResults[0]
	for i := range searchResults {
		if searchResults[i].Length != -1 {
			seedSong = searchResults[i]
		}
	}
	Truncate()
	enqueue(seedSong)
}

func GetPlaylist() Playlist {
	var id int
	var videoid string
	var length int
	var seek int
	playlist := []Song{}
	db := GetDbHandle()
	defer db.Close()
	rows, err := db.Query("SELECT  * from  playlist order by id")
	CheckError(err)
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &videoid, &length, &seek)
		s := Song{
			Id:      id,
			Videoid: videoid,
			Length:  length,
			Seek:    seek,
		}
		CheckError(err)
		playlist = append(playlist, s)
	}
	return playlist
}

func CurrentlyPlaying() Song {
	var id int
	var videoid string
	var length int
	var seek int
	db := GetDbHandle()
	defer db.Close()
	err := db.QueryRow("SELECT * FROM playlist ORDER BY id ASC LIMIT 1").Scan(&id, &videoid, &length, &seek)
	CheckError(err)
	s := Song{
		Id:      id,
		Videoid: videoid,
		Length:  length,
		Seek:    seek,
	}
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

func lastSong() Song {
	var id int
	var videoid string
	var length int
	var seek int
	db := GetDbHandle()
	defer db.Close()
	err := db.QueryRow("SELECT * FROM playlist ORDER BY id DESC LIMIT 1").Scan(&id, &videoid, &length, &seek)
	CheckError(err)
	LastSong := Song{
		Id:      id,
		Videoid: videoid,
		Length:  length,
		Seek:    seek,
	}
	return LastSong
}

func enqueue(s Song) {
	db := GetDbHandle()
	defer db.Close()
	stmt, err := db.Prepare("INSERT INTO playlist (videoid, length, seek) VALUES (?, ?, ?)")
	CheckError(err)
	_, err = stmt.Exec(s.Videoid, s.Length, s.Seek)
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

func autoAdd() {
	ticker := time.NewTicker(time.Second * 5)
	for _ = range ticker.C {
		if Size() == 1 {
			newSong := recommend(lastSong())
			enqueue(newSong)
		}
	}
}

func recommend(s Song) Song {
	recommendations := Recommend(s.Videoid)
	for i := range recommendations {
		if recommendations[i].Length != -1 {
			return recommendations[i]
		}
	}
	return recommendations[0]
}

func Refresh() {
	s := CurrentlyPlaying()
	if s.Seek < s.Length {
		updateSeek(s)
		fmt.Println(s.Videoid, "   ", s.Seek, "          ", GetPlaylist())
	} else {
		remove(s)
		Refresh()
	}
}
