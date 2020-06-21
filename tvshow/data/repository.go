package data

import (
	"github.com/0x113/x-media/tvshow/models"
)

// TVShowRepository contains all methods for operation on TVShow model
type TVShowRepository interface {
	Save(tvShow *models.TVShow) error
	GetByName(name string) (*models.TVShow, error)
	Update(tvShow *models.TVShow) error
	GetAll() ([]*models.TVShow, error)
}
