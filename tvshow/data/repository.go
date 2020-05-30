package data

import "github.com/0x113/x-media/tvshow/models"

// TVShowRepository contains all methods for operation og TVShow model
type TVShowRepository interface {
	Save(tvShow *models.TVShow) error
}
