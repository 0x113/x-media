package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// swagger:response tvShow
// TVShow information
type TVShow struct {
	ID        primitive.ObjectID `bson:"_id" json:"id" validate:"omitempty" example:"507f1f77bcf86cd799439011"`
	Name      string             `bson:"name" json:"name" validate:"required" example:"BoJack Horseman"`
	Language  string             `bson:"language" json:"language" validate:"required" example:"English"`
	Genres    []string           `bson:"genres" json:"genres" validate:"required" example:"Comedy,Drama"`
	Runtime   int                `bson:"runtime" json:"runtime" validate:"required" example:25`
	Premiered string             `bson:"premiered" json:"premiered" validate:"required" example:"2014-08-22"`
	Rating    float32            `bson:"rating" json:"rating" validate:"required" example:"8.1"`
	PosterURL string             `bson:"poster_url" json:"poster_url" validate:"required,url" example:"https://static.tvmaze.com/uploads/images/original_untouched/236/590384.jpg"`
	Summary   string             `bson:"summary" json:"summary" validate:"required" example:"Meet the most beloved sitcom horse of the '90s, 20 years later."`
	DirPath   string             `bson:"dir_path" json:"dir_path" validate:"required" example:"tvshows/BoJack Horseman"`
}
