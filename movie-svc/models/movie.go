package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Movie defines the movie model
type Movie struct {
	ID               primitive.ObjectID `json:"id"`
	TMDbID           int                `json:"tmdb_id" validate:"required"`
	IMDbID           string             `json:"imdb_id"`
	Title            string             `json:"title" validate:"required"`
	Overview         string             `json:"overview" validate:"required"`
	OriginalTitle    string             `json:"original_title" validate:"required"`
	OriginalLanguage string             `json:"original_language" validate:"required"`
	ReleaseDate      string             `json:"release_date" validate:"required"`
	Genres           []string           `json:"genres" validate:"required"`
	Rating           float32            `json:"rating" validate:"required"`
	VoteCount        int                `json:"vote_count" validate:"required"`
	Runtime          int                `json:"runtime"`
	BackdropPath     string             `json:"backdrop_path" validate:"required"`
	PosterPath       string             `json:"poster_path" validate:"required"`
	DirPath          string             `json:"dir_path" validate:"required"`
}
