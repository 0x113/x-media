package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// TVShow information
type TVShow struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Name      string             `bson:"name" json:"name"`
	Language  string             `bson:"language" json:"language"`
	Genres    []string           `bson:"genres" json:"genres"`
	Runtime   uint8              `bson:"runtime" json:"runtime"`
	Premiered string             `bson:"premiered" json:"premiered"`
	Rating    float32            `bson:"rating" json:"rating"`
	PosterURL string             `bson:"poster_url" json:"poster_url"`
	Summary   string             `bson:"summary" json:"summary"`
}
