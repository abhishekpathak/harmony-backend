package musicservice

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

func (s *Song) Score() float64 {
	var score float64
	defer func() {
		if err := recover(); err != nil {
			score = -1.00
		}
	}()
	score = float64(s.Details.Likes-s.Details.Dislikes) * 100.00 / float64(s.Details.Likes)
	return score
}
