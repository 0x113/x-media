package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Movie defines the movie model
type Movie struct {
	ID               primitive.ObjectID `bson:"_id" json:"_id"`
	TMDbID           int                `bson:"tmdb_id" json:"tmdb_id" validate:"required"`
	IMDbID           string             `bson:"imdb_id" json:"imdb_id"`
	Title            string             `bson:"title" json:"title" validate:"required"`
	Overview         string             `bson:"overview" json:"overview" validate:"required"`
	OriginalTitle    string             `bson:"original_title" json:"original_title" validate:"required"`
	OriginalLanguage string             `bson:"original_language" json:"original_language" validate:"required"`
	ReleaseDate      string             `bson:"release_date" json:"release_date" validate:"required"`
	Genres           []string           `bson:"genres" json:"genres" validate:"required"`
	Rating           float32            `bson:"rating" json:"rating" validate:"required"`
	VoteCount        int                `bson:"vote_count" json:"vote_count" validate:"required"`
	Runtime          int                `bson:"runtime" json:"runtime"`
	BackdropPath     string             `bson:"backdrop_path" json:"backdrop_path" validate:"required"`
	PosterPath       string             `bson:"poster_path" json:"poster_path" validate:"required"`
	DirPath          string             `bson:"dir_path" json:"dir_path" validate:"required"`
}
