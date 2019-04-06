package video

// Movie represtents model for movie
type Movie struct {
	MovieID     int64  `json:"movie_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Director    string `json:"director"`
	Genre       string `json:"genre"`
	Duration    string `json:"duration"`
	Rate        string `json:"rate"`
	ReleaseDate string `json:"release_date"`
	FilePath    string `json:"file_path"`
	PosterPath  string `json:"poster_path"`
}
