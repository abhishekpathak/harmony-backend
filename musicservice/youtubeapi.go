package musicservice

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"sort"
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
		searchResults = append(searchResults, CreateSong(item.Id.VideoId))
	}
	return cleanup(searchResults)
}

func getRecommendedResults(videoid string) Playlist {
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
		recommendations = append(recommendations, CreateSong(item.Id.VideoId))
	}
	return cleanup(recommendations)
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
		v.Duration = parseISO8601Duration(string(item.ContentDetails.Duration))
		v.Thumbnail = item.Snippet.Thumbnails.Default.url
		v.Views, _ = strconv.Atoi(item.Statistics.Viewcount)
		v.Likes, _ = strconv.Atoi(item.Statistics.Likecount)
		v.Dislikes, _ = strconv.Atoi(item.Statistics.Dislikecount)
		v.Favourites, _ = strconv.Atoi(item.Statistics.Favouritecount)
		v.Comments, _ = strconv.Atoi(item.Statistics.Commentcount)
	}
	return v
}

func parseISO8601Duration(isoStr string) int {
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

func cleanup(results []Song) []Song {
	var cleanedResults []Song
	for i := range results {
		if results[i].Length != -1 && results[i].Details.Views > 45000 {
			cleanedResults = append(cleanedResults, results[i])
		}
	}
	return cleanedResults
}

func CreateSong(videoid string) Song {
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

func Recommend(s Song) Song {
	var recommendedSong Song
	recommendations := getRecommendedResults(s.Videoid)
	if len(recommendations) < 6 {
		seedQuery := "tum se hi"
		searchResults := Search(seedQuery)
		recommendedSong = searchResults[0]
	} else {
		// sort in the reverse order, so that highest scores come first
		sort.Sort(sort.Reverse(recommendations))
		songindex := rand.Intn(5)
		recommendedSong = recommendations[songindex]
	}
	return recommendedSong
}

func CheckError(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func pprint(songs []Song) string {
	result := "\n"
	for i := range songs {
		result += songs[i].Details.Name + "\t\t" + strconv.FormatFloat(songs[i].Score(), 'f', 2, 64)
		result += "\n"
	}
	return result
}
