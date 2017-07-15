package main

import (
	"fmt"

	"github.com/abhishekpathak/songster/musicservice"
)

func addToLibrary(s musicservice.LibSong, u musicservice.User) bool {
	db := GetDbHandle()
	defer db.Close()
	_, err := db.Exec("insert into library(userid, username, videoid, track, fav) VALUES ($1, $2, $3, $4, $5)", u.Id, u.Name, s.Videoid, s.Track, s.Fav)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func removeFromLibrary(s musicservice.LibSong, u musicservice.User) bool {
	db := GetDbHandle()
	defer db.Close()
	_, err := db.Exec("delete from library where userid = $1 and videoid = $2", u.Id, s.Videoid)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func songExistsInLibrary(userid string, videoid string) bool {
	var size int
	db := GetDbHandle()
	defer db.Close()
	err := db.QueryRow("SELECT count(*) FROM library where userid = $1 and videoid = $2", userid, videoid).Scan(&size)
	if err != nil {
		fmt.Println(err)
		return false
	}
	if size == 0 {
		return false
	}
	return true
}

func randomSongFromLibrary(userid string) musicservice.LibSong {
	libSongObj := musicservice.LibSong{}
	db := GetDbHandle()
	defer db.Close()
	query := fmt.Sprintf("SELECT videoid, track, fav, source FROM library where userid = '%s' and last_played not in(select max(last_played) from library where userid = '%s') order by random() limit 1", userid, userid)
	err := db.QueryRow(query).Scan(&libSongObj.Videoid, &libSongObj.Track, &libSongObj.Fav, &libSongObj.Source)
	if err != nil {
		fmt.Println(err)
	}
	return libSongObj
}

func updateLastPlayedTimestamp(userid string, videoid string) bool {
	db := GetDbHandle()
	defer db.Close()
	_, err := db.Exec("update library set last_played = current_timestamp where userid = $1 and videoid = $2", userid, videoid)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}
