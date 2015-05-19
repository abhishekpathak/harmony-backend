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
	searchUrl += "&part=snippet&fields=items(id%2Csnippet)&videoSyndicated=true&type=video&videoDuration=medium"
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
		searchResults = append(searchResults, createSong(item.Id.VideoId, item.Snippet.Title))
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
	recommendUrl += "&part=snippet&fields=items(id%2Csnippet)&videoSyndicated=true&type=video&videoDuration=medium"
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
		recommendations = append(recommendations, createSong(item.Id.VideoId, item.Snippet.Title))
	}
	return recommendations
}

func getDuration(videoid string) int {
	type Item struct {
		Duration        string
		Dimension       string
		definition      string
		caption         string
		licensedContent string
	}

	type ContentDetails struct {
		ContentDetails Item
	}

	type Resp struct {
		Items []ContentDetails
	}
	recommendUrl := "https://www.googleapis.com/youtube/v3/videos?part=snippet%2CcontentDetails"
	recommendUrl += fmt.Sprintf("&id=%s", videoid)
	recommendUrl += "&fields=items%2FcontentDetails"
	recommendUrl += fmt.Sprintf("&key=%s", API_KEY)

	response, err := http.Get(recommendUrl)
	CheckError(err)
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	CheckError(err)

	resp := Resp{}
	err = json.Unmarshal([]byte(contents), &resp)
	CheckError(err)
	ISO8601Duration := string(resp.Items[0].ContentDetails.Duration)
	return ParseISO8601Duration(ISO8601Duration)
}

func ParseISO8601Duration(isoStr string) int {
	//PT6M11S
	//PT41M44S
	//PT1H18M27S
	isoStr = strings.Replace(isoStr, "PT", "", 1)
	isoStr = strings.Replace(isoStr, "H", ",", 1)
	isoStr = strings.Replace(isoStr, "M", ",", 1)
	isoStr = strings.Replace(isoStr, "S", "", 1)
	timeSlice := strings.Split(isoStr, ",")
	if len(timeSlice) > 2 {
		return -1
	}
	minutes, err := strconv.Atoi(timeSlice[0])
	CheckError(err)
	seconds, err := strconv.Atoi(timeSlice[1])
	CheckError(err)
	duration := minutes*60 + seconds
	if duration < 120 || duration > 600 {
		duration = -1
	}
	return duration
}
