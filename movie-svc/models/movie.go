package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Movie defines the movie model
type Movie struct {
	ID               primitive.ObjectID `bson:"_id" json:"_id" example:"507f1f77bcf86cd799439011"`
	TMDbID           int                `bson:"tmdb_id" json:"tmdb_id" validate:"required" example:"949"`
	IMDbID           string             `bson:"imdb_id" json:"imdb_id" example:"tt0113277"`
	Title            string             `bson:"title" json:"title" validate:"required" example:"Heat"`
	Overview         string             `bson:"overview" json:"overview" validate:"required" example:"Obsessive master thief, Neil McCauley leads a top-notch crew on various daring heists throughout Los Angeles while determined detective, Vincent Hanna pursues him without rest. Each man recognizes and respects the ability and the dedication of the other even though they are aware their cat-and-mouse game may end in violence."`
	OriginalTitle    string             `bson:"original_title" json:"original_title" validate:"required" example:"Heat"`
	OriginalLanguage string             `bson:"original_language" json:"original_language" validate:"required" example:"en"`
	ReleaseDate      string             `bson:"release_date" json:"release_date" validate:"required" example:"1995-12-15"`
	Genres           []string           `bson:"genres" json:"genres" validate:"required" example:"Action,Crime,Drama,Thriller"`
	Rating           float32            `bson:"rating" json:"rating" validate:"required" example:"7.9"`
	VoteCount        int                `bson:"vote_count" json:"vote_count" validate:"required" example:"420"`
	Runtime          int                `bson:"runtime" json:"runtime" example:"170"`
	BackdropPath     string             `bson:"backdrop_path" json:"backdrop_path" validate:"required" example:"/rfEXNlql4CafRmtgp2VFQrBC4sh.jpg"`
	PosterPath       string             `bson:"poster_path" json:"poster_path" validate:"required" example:"/rrBuGu0Pjq7Y2BWSI6teGfZzviY.jpg"`
	DirPath          string             `bson:"dir_path" json:"dir_path" validate:"required" example:"/home/0x113/Movies/Heat.1995.mp4"`
}
