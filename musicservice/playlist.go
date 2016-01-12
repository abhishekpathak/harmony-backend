package musicservice

type Playlist []Song

func (s Playlist) Len() int {
	return len(s)
}

func (s Playlist) Less(i, j int) bool {
	if s[i].Score() <= s[j].Score() {
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
