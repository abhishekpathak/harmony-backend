package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const API_KEY = "AIzaSyCbfxhEDNKXXPFbmjttsqFvGHxjvTlfVxg"

type SongInfo struct {
	Name       string
	Duration   int
	Thumbnail  string
	Views      int
	Likes      int
	Dislikes   int
	Favourites int
	Comments   int
}

func (v *SongInfo) init() SongInfo {
	return SongInfo{
		Name:       "not found",
		Duration:   -1,
		Thumbnail:  "not found",
		Views:      -1,
		Likes:      -1,
		Dislikes:   -1,
		Favourites: -1,
		Comments:   -1,
	}
}

func Search(query string) []Song {
	type Id struct {
		Kind    string
		VideoId string
	}

	type Url struct {
		url string
	}

	type Thumbnail struct {
		Default Url
		Medium  Url
		High    Url
	}

	type Snippet struct {
		PublishedAt          string
		ChannelId            string
		Title                string
		Description          string
		Thumbnails           Thumbnail
		ChannelTitle         string
		LiveBroadcastContent string
	}

	type Item struct {
		Id      Id
		Snippet Snippet
	}

	type Resp struct {
		Items []Item
	}

	searchUrl := fmt.Sprintf("https://www.googleapis.com/youtube/v3/search?videoEmbeddable=true&q=%s", url.QueryEscape(query))
	searchUrl += "&part=snippet&fields=items(id%2Csnippet)&type=video&maxResults=5"
	searchUrl += fmt.Sprintf("&key=%s", API_KEY)

	response, err := http.Get(searchUrl)
	CheckError(err)
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	CheckError(err)

	resp := Resp{}
	err = json.Unmarshal([]byte(contents), &resp)
	CheckError(err)
	searchResults := []Song{}

	for _, item := range resp.Items {
		searchResults = append(searchResults, createSong(item.Id.VideoId))
	}
	return searchResults
}

func Recommend(videoid string) []Song {
	type Id struct {
		Kind    string
		VideoId string
	}

	type Url struct {
		url string
	}

	type Thumbnail struct {
		Default Url
		Medium  Url
		High    Url
	}

	type Snippet struct {
		PublishedAt          string
		ChannelId            string
		Title                string
		Description          string
		Thumbnails           Thumbnail
		ChannelTitle         string
		LiveBroadcastContent string
	}

	type Item struct {
		Id      Id
		Snippet Snippet
	}

	type Resp struct {
		Items []Item
	}

	recommendUrl := fmt.Sprintf("https://www.googleapis.com/youtube/v3/search?videoEmbeddable=true&relatedToVideoId=%s", videoid)
	recommendUrl += "&part=snippet&fields=items(id%2Csnippet)&type=video&maxResults=20"
	recommendUrl += fmt.Sprintf("&key=%s", API_KEY)

	response, err := http.Get(recommendUrl)
	CheckError(err)
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	CheckError(err)

	resp := Resp{}
	err = json.Unmarshal([]byte(contents), &resp)
	CheckError(err)
	recommendations := []Song{}

	for _, item := range resp.Items {
		recommendations = append(recommendations, createSong(item.Id.VideoId))
	}
	return recommendations
}

func GetInfo(videoid string) SongInfo {
	type Url struct {
		url string
	}

	type Thumbnail struct {
		Default Url
		Medium  Url
		High    Url
	}

	type Localized struct {
		Title       string
		description string
	}

	type Snippet struct {
		PublishedAt          string
		ChannelId            string
		Title                string
		Description          string
		Thumbnails           Thumbnail
		ChannelTitle         string
		CategoryId           string
		LiveBroadcastContent string
		Localized            Localized
	}

	type ContentDetails struct {
		Duration        string
		Dimension       string
		definition      string
		caption         string
		licensedContent string
	}

	type Statistics struct {
		Viewcount      string
		Likecount      string
		Dislikecount   string
		Favouritecount string
		Commentcount   string
	}

	type Item struct {
		Id             string
		Snippet        Snippet
		ContentDetails ContentDetails
		Statistics     Statistics
	}

	type Resp struct {
		Items []Item
	}

	infoUrl := "https://www.googleapis.com/youtube/v3/videos?part=snippet%2CcontentDetails%2Cstatistics"
	infoUrl += fmt.Sprintf("&id=%s", videoid)
	infoUrl += "&fields=items(contentDetails%2Cid%2Csnippet%2Cstatistics%2Csuggestions)"
	infoUrl += fmt.Sprintf("&key=%s", API_KEY)

	response, err := http.Get(infoUrl)
	CheckError(err)
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	CheckError(err)

	resp := Resp{}
	err = json.Unmarshal([]byte(contents), &resp)
	CheckError(err)
	var v = SongInfo{}
	v = v.init()
	if len(resp.Items) < 1 {
		return v
	}
	item := resp.Items[0]

	if item.Snippet.CategoryId == "10" {
		v.Name = item.Snippet.Title
		v.Duration = ParseISO8601Duration(string(item.ContentDetails.Duration))
		v.Thumbnail = item.Snippet.Thumbnails.Default.url
		v.Views, _ = strconv.Atoi(item.Statistics.Viewcount)
		v.Likes, _ = strconv.Atoi(item.Statistics.Likecount)
		v.Dislikes, _ = strconv.Atoi(item.Statistics.Dislikecount)
		v.Favourites, _ = strconv.Atoi(item.Statistics.Favouritecount)
		v.Comments, _ = strconv.Atoi(item.Statistics.Commentcount)
	}
	return v
}

func ParseISO8601Duration(isoStr string) int {
	//PT6M11S
	//PT41M44S
	//PT1H18M27S
	//PT15M
	isoStr = strings.Replace(isoStr, "PT", "", 1)
	isoStr = strings.Replace(isoStr, "H", ",", 1)
	isoStr = strings.Replace(isoStr, "M", ",", 1)
	isoStr = strings.Replace(isoStr, "S", "", 1)
	timeSlice := strings.Split(isoStr, ",")
	if len(timeSlice) != 2 {
		return -1
	}
	minutes, err := strconv.Atoi(timeSlice[0])
	if err != nil {
		return -1
	}
	seconds, err := strconv.Atoi(timeSlice[1])
	if err != nil {
		return -1
	}
	duration := minutes*60 + seconds
	if duration < 120 || duration > 600 {
		duration = -1
	}
	return duration
}
