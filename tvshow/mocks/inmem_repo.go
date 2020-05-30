package mocks

import (
	"errors"

	"github.com/0x113/x-media/tvshow/models"
)

// ITVShowRepo represents in-memory tv show repository
type ITVShowRepo struct {
	tvShows map[string]*models.TVShow
}

// NewInmemTVShowRepository creates new ITVShowRepo
func NewInmemTVShowRepository() *ITVShowRepo {
	var tvShows = map[string]*models.TVShow{}
	return &ITVShowRepo{tvShows}
}

// Save tv show in memory
func (r *ITVShowRepo) Save(tvShow *models.TVShow) error {
	if _, ok := r.tvShows[tvShow.ID.Hex()]; ok {
		return errors.New("TV Show already exists")
	}
	r.tvShows[tvShow.ID.Hex()] = tvShow
	return nil
}
