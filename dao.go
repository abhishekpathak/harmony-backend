package main

import "github.com/HarmonyProject/songster/musicservice"

func getLastPlaying(userId string) musicservice.LibSong {
	var l musicservice.LibSong
	db := GetDbHandle()
	defer db.Close()
	err := db.QueryRow("SELECT videoid, artist, track, rating, fav from library where userid = ? order by last_played desc limit 1", userId).Scan(&l.Videoid, &l.Artist, &l.Track, &l.Rating, &l.Fav)
	CheckError(err)
	return l
}

func addToLibrary(s musicservice.LibSong, u musicservice.User) bool {
	db := GetDbHandle()
	defer db.Close()
	stmt, err := db.Prepare("insert into library(userid, username, videoid, artist, track, rating, fav) VALUES (?, ?, ?, ?, ?, ?, ?)")
	CheckError(err)
	_, err = stmt.Exec(u.Id, u.Name, s.Videoid, s.Artist, s.Track, s.Rating, s.Fav)
	CheckError(err)
	return true
}

func removeFromLibrary(s musicservice.LibSong, u musicservice.User) bool {
	db := GetDbHandle()
	defer db.Close()
	stmt, err := db.Prepare("delete from library where userid = ? and videoid = ?")
	CheckError(err)
	_, err = stmt.Exec(u.Id, s.Videoid)
	CheckError(err)
	return true
}

func songExistsInLibrary(userid string, videoid string) bool {
	var size int
	db := GetDbHandle()
	defer db.Close()
	err := db.QueryRow("SELECT count(*) FROM library where userid = ? and videoid = ?", userid, videoid).Scan(&size)
	CheckError(err)
	return true
}

func favSongFromLibrary(userid string) musicservice.LibSong {
	libSongObj := musicservice.LibSong{}
	db := GetDbHandle()
	defer db.Close()
	err := db.QueryRow("SELECT videoid, artist, track, rating, fav FROM library where userid = ? and fav = true and last_played not in(select max(last_played) from library where userid = ?) order by rand() limit 1", userid, userid).Scan(&libSongObj.Videoid, &libSongObj.Artist, &libSongObj.Track, &libSongObj.Rating, &libSongObj.Fav)
	CheckError(err)
	return libSongObj
}

func randomSongFromLibrary(userid string) musicservice.LibSong {
	libSongObj := musicservice.LibSong{}
	db := GetDbHandle()
	defer db.Close()
	err := db.QueryRow("SELECT videoid, artist, track, rating, fav FROM library where userid = ? and last_played not in(select max(last_played) from library where userid = ?) order by rand() limit 1", userid, userid).Scan(&libSongObj.Videoid, &libSongObj.Artist, &libSongObj.Track, &libSongObj.Rating, &libSongObj.Fav)
	CheckError(err)
	return libSongObj
}

func updateLastPlayedTimestamp(userid string, videoid string) bool {
	db := GetDbHandle()
	defer db.Close()
	stmt, err := db.Prepare("update library set last_played = current_timestamp where userid = ? and videoid = ?")
	CheckError(err)
	_, err = stmt.Exec(userid, videoid)
	CheckError(err)
	return true
}
