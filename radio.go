package main

import (
	"fmt"
	_ "io/ioutil"
	"strconv"
	"time"
)

var i = 0

type song struct {
	name   string
	length int
	seek   int
}

type playlist struct {
	songs []song
}

func getRecommendedSong(currentSong song) song {
	i = i + 1
	recommendedSong := song{
		name:   "Song" + strconv.Itoa(i),
		length: 5,
		seek:   0,
	}
	time.Sleep(3 * time.Second)
	return recommendedSong
}

func (p *playlist) seed(ss song) {
	p.songs = append(p.songs, ss)
}

func (p *playlist) getCurrentlyPlaying() song {
	return p.songs[0]
}

func (p *playlist) setCurrentlyPlaying(s song) {
	p.songs[0] = s
}

func (p *playlist) getLastSong() song {
	return p.songs[len(p.songs)-1]
}

func (p *playlist) enqueue(newSong song) {
	p.songs = append(p.songs, newSong)
}

func (p *playlist) removeFromTop() {
	p.songs = p.songs[1:]
}

func (p *playlist) autoAdd() {
	if len(p.songs) == 1 {
		nextSong := getRecommendedSong(p.getLastSong())
		p.enqueue(nextSong)
	}
}

func (p *playlist) refresh() {
	s := p.getCurrentlyPlaying()
	if s.seek < s.length {
		s.seek++
		p.setCurrentlyPlaying(s)
		fmt.Println(s.name, "   ", s.seek, "          ", p.songs)
	} else {
		p.removeFromTop()
		p.refresh()
	}
}

func main() {
	myplaylist := &playlist{}
	seedSong := song{
		name:   "Song" + strconv.Itoa(i),
		length: 5,
		seek:   0,
	}
	myplaylist.seed(seedSong)
	ticker := time.NewTicker(time.Second)
	for _ = range ticker.C {
		go myplaylist.autoAdd()
		myplaylist.refresh()
	}
}
