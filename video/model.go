package video

// Movie represtents model for movie
type Movie struct {
	MovieID     int64   `json:"movie_id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Director    string  `json:"director"`
	Genre       string  `json:"genre"`
	Duration    string  `json:"duration"`
	Rate        float64 `json:"rate"`
	ReleaseDate string  `json:"release_date"`
	PosterPath  string  `json:"poster_path"`
}

type TVSeries struct {
	TvSeriesID      int64   `json:"tv_series_id"`
	Title           string  `json:"title"`
	Description     string  `json:"description"`
	Director        string  `json:"director"`
	Genre           string  `json:"genre"`
	EpisodeDuration string  `json:"episode_duration"`
	Rate            float64 `json:"rate"`
	ReleaseDate     string  `json:"release_date"`
	PosterPath      string  `json:"poster_path"`
}
