package main

import (
	"fmt"

	"github.com/HarmonyProject/songster/musicservice"
)

func getLastPlaying(userId string) musicservice.LibSong {
	var l musicservice.LibSong
	db := GetDbHandle()
	defer db.Close()
	err := db.QueryRow("SELECT videoid, artist, track, rating, fav from library where userid = ? order by last_played desc limit 1", userId).Scan(&l.Videoid, &l.Artist, &l.Track, &l.Rating, &l.Fav)
	if err != nil {
		fmt.Println(err)
		fmt.Println("unable to find the last played song. Picking a random song.")
		return randomSongFromLibrary(userId)
	}
	return l
}

func addToLibrary(s musicservice.LibSong, u musicservice.User) bool {
	db := GetDbHandle()
	defer db.Close()
	stmt, err := db.Prepare("insert into library(userid, username, videoid, artist, track, rating, fav) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		fmt.Println(err)
		return false
	}
	_, err = stmt.Exec(u.Id, u.Name, s.Videoid, s.Artist, s.Track, s.Rating, s.Fav)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func removeFromLibrary(s musicservice.LibSong, u musicservice.User) bool {
	db := GetDbHandle()
	defer db.Close()
	stmt, err := db.Prepare("delete from library where userid = ? and videoid = ?")
	if err != nil {
		fmt.Println(err)
		return false
	}
	_, err = stmt.Exec(u.Id, s.Videoid)
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
	err := db.QueryRow("SELECT count(*) FROM library where userid = ? and videoid = ?", userid, videoid).Scan(&size)
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
	err := db.QueryRow("SELECT videoid, artist, track, rating, fav FROM library where userid = ? and last_played not in(select max(last_played) from library where userid = ?) order by rand() limit 1", userid, userid).Scan(&libSongObj.Videoid, &libSongObj.Artist, &libSongObj.Track, &libSongObj.Rating, &libSongObj.Fav)
	if err != nil {
		fmt.Println(err)
	}
	return libSongObj
}

func updateLastPlayedTimestamp(userid string, videoid string) bool {
	db := GetDbHandle()
	defer db.Close()
	stmt, err := db.Prepare("update library set last_played = current_timestamp where userid = ? and videoid = ?")
	if err != nil {
		fmt.Println(err)
		return false
	}
	_, err = stmt.Exec(userid, videoid)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}
