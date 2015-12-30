package main

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"
)

// Song will have all data that Melody requires, nothing more.
type Song struct {
	Id        int      `json:"id"`
	Videoid   string   `json:"videoid"`
	Name      string   `json:"name"`
	Length    int      `json:"length"`
	Seek      int      `json:"seek"`
	AddedBy   string   `json:"added_by"`
	Thumbnail string   `json:"thumbnail"`
	Details   SongInfo `json:"details"`
}

type Playlist []Song

func (s Playlist) Len() int {
	return len(s)
}

func (s Playlist) Less(i, j int) bool {
	if s[i].score() <= s[j].score() {
		return true
	} else {
		return false
	}
}

func (s Playlist) Swap(i, j int) {
	temp := s[i]
	s[i] = s[j]
	s[j] = temp
}

func (s *Song) score() int {
	return s.Details.Views/10000 + (s.Details.Likes-s.Details.Dislikes)/10 + s.Details.Favourites*10 + s.Details.Comments
}

func getSong(id int) Song {
	var videoid string
	var name string
	var length int
	var seek int
	var addedBy string
	var thumbnail string
	var songInfo SongInfo

	query := "select * from playlist where id = ?"
	db := GetDbHandle()
	defer db.Close()
	err := db.QueryRow(query, id).Scan(&id, &videoid, &name, &length, &seek, &addedBy, &thumbnail)
	CheckError(err)
	query = "select name, duration, thumbnail, views, likes, dislikes, favourites, comments from song_details where videoid = ?"
	err = db.QueryRow(query, videoid).Scan(&songInfo.Name, &songInfo.Duration, &songInfo.Thumbnail, &songInfo.Views, &songInfo.Likes, &songInfo.Dislikes, &songInfo.Favourites, &songInfo.Comments)
	CheckError(err)

	return Song{
		Id:        id,
		Videoid:   videoid,
		Name:      name,
		Length:    length,
		Seek:      seek,
		AddedBy:   addedBy,
		Thumbnail: thumbnail,
		Details:   songInfo,
	}
}

func createSong(videoid string) Song {
	details := GetInfo(videoid)
	return Song{
		Id:        -1,
		Videoid:   videoid,
		Name:      details.Name,
		Length:    details.Duration,
		Seek:      -5,
		AddedBy:   "system",
		Thumbnail: details.Thumbnail,
		Details:   details,
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

func cleanup(results []Song) Playlist {
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

func GetPlaylist() []string {
	var name string
	playlist := make([]string, 1)
	db := GetDbHandle()
	defer db.Close()
	rows, err := db.Query("SELECT name from  playlist order by id")
	CheckError(err)
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&name)
		//s := getSong(id)
		CheckError(err)
		playlist = append(playlist, name)
	}
	return playlist
}

func CurrentlyPlaying() Song {
	var id int
	db := GetDbHandle()
	defer db.Close()
	err := db.QueryRow("SELECT id FROM playlist ORDER BY id ASC LIMIT 1").Scan(&id)
	CheckError(err)
	s := getSong(id)
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
	db := GetDbHandle()
	defer db.Close()
	err := db.QueryRow("SELECT id FROM playlist ORDER BY id DESC LIMIT 1").Scan(&id)
	CheckError(err)
	lastSong := getSong(id)
	return lastSong
}

func enqueue(s Song, agent string) {
	db := GetDbHandle()
	defer db.Close()
	stmt, err := db.Prepare("INSERT INTO playlist (videoid, name, length, seek, added_by) VALUES (?, ?, ?, ?, ?)")
	CheckError(err)
	_, err = stmt.Exec(s.Videoid, s.Name, s.Length, s.Seek, agent)
	CheckError(err)
	stmt, err = db.Prepare("REPLACE INTO song_details(videoid, name, duration, thumbnail, views, likes, dislikes,  favourites, comments) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)")
	CheckError(err)
	_, err = stmt.Exec(s.Videoid, s.Details.Name, s.Details.Duration, s.Details.Thumbnail, s.Details.Views, s.Details.Likes, s.Details.Dislikes, s.Details.Favourites, s.Details.Comments)
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
		song := createSong(query)
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
		sort.Sort(recommendations)
		songindex := rand.Intn(5)
		recommendedSong = recommendations[songindex]
	}
	return recommendedSong
}

func Refresh() {
	s := CurrentlyPlaying()
	if s.Seek < s.Length {
		updateSeek(s)
		fmt.Println(s.Seek, "          ", GetPlaylist())
	} else {
		remove(s)
		//go PostToSlack("#nowplaying " + CurrentlyPlaying().Name)
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
