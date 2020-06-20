package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// TVShow information
type TVShow struct {
	ID        primitive.ObjectID `bson:"_id" json:"id" validate:"omitempty"`
	Name      string             `bson:"name" json:"name" validate:"required"`
	Language  string             `bson:"language" json:"language" validate:"required"`
	Genres    []string           `bson:"genres" json:"genres" validate:"required"`
	Runtime   int                `bson:"runtime" json:"runtime" validate:"required"`
	Premiered string             `bson:"premiered" json:"premiered" validate:"required"`
	Rating    float32            `bson:"rating" json:"rating" validate:"required"`
	PosterURL string             `bson:"poster_url" json:"poster_url" validate:"required,url"`
	Summary   string             `bson:"summary" json:"summary" validate:"required"`
	DirPath   string             `bson:"dir_path" json:"dir_path" validate:"required"`
}
