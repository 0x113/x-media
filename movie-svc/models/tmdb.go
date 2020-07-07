package models

// TMDbQueryResponse represents response for the https://api.themoviedb.org/3/search/movie?api_key={api_key}&query={title}&language={lang}
type TMDbQueryResponse struct {
	Page         int               `json:"page"`
	TotalResults int               `json:"total_results"`
	TotalPages   int               `json:"total_pages"`
	Results      []*TMDbQueryMovie `json:"results"`
}

// TMDbQueryMovie represents one result model from the https://api.themoviedb.org/3/search/movie?api_key={api_key}&query={title}&language={lang}
type TMDbQueryMovie struct {
	Popularity       float64 `json:"popularity"`
	ID               int     `json:"id"`
	Video            bool    `json:"video"`
	VoteCount        int     `json:"vote_count"`
	VoteAverage      float64 `json:"vote_average"`
	Title            string  `json:"title"`
	ReleaseDate      string  `json:"release_date"`
	OriginalLanguage string  `json:"original_language"`
	OriginalTitle    string  `json:"original_title"`
	GenreIds         []int   `json:"genre_ids"`
	BackdropPath     string  `json:"backdrop_path"`
	Adult            bool    `json:"adult"`
	Overview         string  `json:"overview"`
	PosterPath       string  `json:"poster_path"`
}
