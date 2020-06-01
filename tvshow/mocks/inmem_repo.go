package mocks

import (
	"errors"

	"github.com/0x113/x-media/tvshow/models"
)

// MockTVShowRepository represents in-memory tv show repository
type MockTVShowRepository struct {
	tvShows map[string]*models.TVShow
}

// NewInmemTVShowRepository creates new MockTVShowRepository
func NewMockTVShowRepository() *MockTVShowRepository {
	var tvShows = map[string]*models.TVShow{}
	return &MockTVShowRepository{tvShows}
}

// Save tv show in memory
func (r *MockTVShowRepository) Save(tvShow *models.TVShow) error {
	if _, ok := r.tvShows[tvShow.ID.Hex()]; ok {
		return errors.New("TV Show already exists")
	}
	r.tvShows[tvShow.ID.Hex()] = tvShow
	return nil
}
