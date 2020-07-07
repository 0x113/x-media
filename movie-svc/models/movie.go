package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Movie defines the movie model
type Movie struct {
	ID               primitive.ObjectID `json:"id"`
	TMDbID           int                `json:"tm_db_id"`
	Title            string             `json:"title"`
	OriginalTitle    string             `json:"original_title"`
	OriginalLanguage string             `json:"original_language"`
	ReleaseDate      string             `json:"release_date"`
	Genres           []string           `json:"genres"`
	Rating           float32            `json:"rating"`
	Runtime          int                `json:"runtime"`
	BackdropPath     string             `json:"backdrop_path"`
	PosterPath       string             `json:"poster_path"`
	DirPath          string             `json:"dir_path"`
}
