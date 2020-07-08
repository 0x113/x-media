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
	Popularity       float32 `json:"popularity"`
	ID               int     `json:"id"`
	Video            bool    `json:"video"`
	VoteCount        int     `json:"vote_count"`
	VoteAverage      float32 `json:"vote_average"`
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

// TMDbGenre defines one genre from the https://api.themoviedb.org/3/genre/movie/list?api_key={api_key}&language={lang}
type TMDbGenre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// TMDbMovie defines the response from https://api.themoviedb.org/3/movie/949?api_key={api_key}&language={lang}
type TMDbMovie struct {
	BackdropPath        string       `json:"backdrop_path"`
	Genres              []*TMDbGenre `json:"genres"`
	ID                  int          `json:"id"`
	IMDbID              string       `json:"imdb_id"`
	OriginalLanguage    string       `json:"original_language"`
	OriginalTitle       string       `json:"original_title"`
	Overview            string       `json:"overview"`
	PosterPath          string       `json:"poster_path"`
	ProductionCountries []struct {
		Iso31661 string `json:"iso_3166_1"`
		Name     string `json:"name"`
	} `json:"production_countries"`
	ReleaseDate string  `json:"release_date"`
	Runtime     int     `json:"runtime"`
	Title       string  `json:"title"`
	Video       bool    `json:"video"`
	VoteAverage float32 `json:"vote_average"`
	VoteCount   int     `json:"vote_count"`
}
