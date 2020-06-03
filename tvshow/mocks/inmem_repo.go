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
	tvShows["BoJack Horseman"] = &models.TVShow{
		Name:      "BoJack Horseman",
		Language:  "English",
		Genres:    []string{"Comedy", "Drama"},
		Runtime:   25,
		Premiered: "2014-08-22",
		Rating:    8.1,
		PosterURL: "https://static.tvmaze.com/uploads/images/original_untouched/236/590384.jpg",
		Summary:   "Meet the most beloved sitcom horse of the '90s, 20 years later.",
	}
	return &MockTVShowRepository{tvShows}
}

// Save tv show in memory
func (r *MockTVShowRepository) Save(tvShow *models.TVShow) error {
	if _, ok := r.tvShows[tvShow.Name]; ok {
		return errors.New("TV Show already exists")
	}
	r.tvShows[tvShow.Name] = tvShow
	return nil
}
