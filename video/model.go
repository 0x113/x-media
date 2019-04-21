package video

// Movie represents model for movie
type Movie struct {
	MovieID     int64   `json:"movie_id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Director    string  `json:"director"`
	Genre       string  `json:"genre"`
	Duration    string  `json:"duration"`
	Rate        float64 `json:"rate"`
	ReleaseDate string  `json:"release_date"`
	FileName    string  `json:"file_name"`
	PosterPath  string  `json:"poster_path"`
}

// TvSeries represents model for tv series
type TVSeries struct {
	SeriesID        int64   `json:"series_id"`
	Title           string  `json:"title"`
	Description     string  `json:"description"`
	Director        string  `json:"director"`
	Genre           string  `json:"genre"`
	EpisodeDuration string  `json:"episode_duration"`
	Rate            float64 `json:"rate"`
	ReleaseDate     string  `json:"release_date"`
	DirName         string  `json:"dir_name"`
	PosterPath      string  `json:"poster_path"`
}

// Season represents model for one season of tv series
type Season struct {
	Name     string   `json:"name"`
	Episodes []string `json:"episodes"`
}
