package musicservice

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
