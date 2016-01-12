package main

import (
	"github.com/HarmonyProject/songster/musicservice"
)

func addToPlaylist(s musicservice.Song, agent string) {
	db := GetDbHandle()
	defer db.Close()
	stmt, err := db.Prepare("INSERT INTO playlist (videoid, name, length, seek, added_by) VALUES (?, ?, ?, ?, ?)")
	CheckError(err)
	_, err = stmt.Exec(s.Videoid, s.Name, s.Length, s.Seek, agent)
	CheckError(err)
}

func getSong(id int) musicservice.Song {
	/*
		var videoid string
		var name string
		var length int
		var seek int
		var addedBy string
		var thumbnail string
	*/
	var song = musicservice.Song{}
	query := "select * from playlist where id = ?"
	db := GetDbHandle()
	defer db.Close()
	err := db.QueryRow(query, id).Scan(&song.Id, &song.Videoid, &song.Name, &song.Length, &song.Seek, &song.AddedBy, &song.Thumbnail)
	CheckError(err)
	song.Details = getSongDetails(song.Videoid)
	return song
}

func getSongDetails(videoid string) musicservice.SongInfo {
	var songInfo = musicservice.SongInfo{}
	query := "select name, duration, thumbnail, views, likes, dislikes, favourites, comments from song_details where videoid = ?"
	db := GetDbHandle()
	defer db.Close()
	err := db.QueryRow(query, videoid).Scan(&songInfo.Name, &songInfo.Duration, &songInfo.Thumbnail, &songInfo.Views, &songInfo.Likes, &songInfo.Dislikes, &songInfo.Favourites, &songInfo.Comments)
	CheckError(err)
	return songInfo
}

func firstSongId() int {
	var id int
	db := GetDbHandle()
	defer db.Close()
	err := db.QueryRow("SELECT id FROM playlist ORDER BY id ASC LIMIT 1").Scan(&id)
	CheckError(err)
	return id
}

func lastSongId() int {
	var id int
	db := GetDbHandle()
	defer db.Close()
	err := db.QueryRow("SELECT id FROM playlist ORDER BY id DESC LIMIT 1").Scan(&id)
	CheckError(err)
	return id
}

func clearPlaylist() {
	db := GetDbHandle()
	defer db.Close()
	stmt, err := db.Prepare("TRUNCATE playlist")
	CheckError(err)
	_, err = stmt.Exec()
	CheckError(err)
}

func removeFromPlaylist(id int) {
	db := GetDbHandle()
	defer db.Close()
	stmt, err := db.Prepare("DELETE FROM playlist WHERE id = ?")
	CheckError(err)
	_, err = stmt.Exec(id)
	CheckError(err)

}

func updateSeek(id int) {
	db := GetDbHandle()
	defer db.Close()
	stmt, err := db.Prepare("UPDATE playlist SET seek = seek + 1 WHERE id = ?")
	CheckError(err)
	_, err = stmt.Exec(id)
	CheckError(err)
}

func currentPlaylistIds() []int {
	var id int
	var ids = make([]int, 1, 2)
	db := GetDbHandle()
	defer db.Close()
	rows, err := db.Query("SELECT id from  playlist order by id")
	CheckError(err)
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id)
		if err == nil {
			ids = append(ids, id)
		}
	}
	return ids
}

func playlistSize() int {
	var size int
	db := GetDbHandle()
	defer db.Close()
	err := db.QueryRow("SELECT count(*) FROM playlist").Scan(&size)
	CheckError(err)
	return size
}

func UpdateSongdetails(s musicservice.Song) {
	db := GetDbHandle()
	defer db.Close()
	stmt, err := db.Prepare("REPLACE INTO song_details(videoid, name, duration, thumbnail, views, likes, dislikes,  favourites, comments, score) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	CheckError(err)
	_, err = stmt.Exec(s.Videoid, s.Details.Name, s.Details.Duration, s.Details.Thumbnail, s.Details.Views, s.Details.Likes, s.Details.Dislikes, s.Details.Favourites, s.Details.Comments, s.Score())
	CheckError(err)
}
