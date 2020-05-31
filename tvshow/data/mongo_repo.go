package data

import (
	"context"

	"github.com/0x113/x-media/tvshow/databases"
	"github.com/0x113/x-media/tvshow/models"
)

const (
	collectionName = "tvshows"
)

// tvShowRepository manages the TVShow CRUD
type tvShowRepository struct{}

// NewMongoTVShowRepository returns new instance of TVShowRepository
func NewMongoTVShowRepository() TVShowRepository {
	return &tvShowRepository{}
}

// Save TVShow to the database
func (r *tvShowRepository) Save(tvShow *models.TVShow) error {
	// get session
	sessionCopy := databases.Database.Session
	defer sessionCopy.EndSession(context.TODO())

	collection := sessionCopy.Client().Database(databases.Database.DbName).Collection(collectionName)

	_, err := collection.InsertOne(context.TODO(), tvShow)
	if err != nil {
		return err
	}

	return nil
}
