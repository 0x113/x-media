package service_test

import (
	"io/ioutil"
	"testing"

	"github.com/0x113/x-media/tvshow/mocks"
	"github.com/0x113/x-media/tvshow/models"
	"github.com/0x113/x-media/tvshow/service"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestSuccessSave(t *testing.T) {
	tvShowRepo := mocks.NewInmemTVShowRepository()
	tvShowService := service.NewTVShowService(tvShowRepo)
	log.SetOutput(ioutil.Discard) // disable logrus for tests

	tvShow := &models.TVShow{
		Name:      "The Office",
		Language:  "English",
		Genres:    []string{"Comedy"},
		Runtime:   30,
		Premiered: "2005-03-24",
		Rating:    8.5,
		PosterURL: "https://static.tvmaze.com/uploads/images/original_untouched/85/213184.jpg",
		Summary:   "Steve Carell stars in The Office, a fresh and funny mockumentary-style glimpse into the daily interactions of the eccentric workers at the Dunder Mifflin paper supply company. Based on the smash-hit British series of the same name and adapted for American Television by Greg Daniels, this fast-paced comedy parodies contemporary American water-cooler culture. Earnest but clueless regional manager Michael Scott believes himself to be an exceptional boss and mentor, but actually receives more eye-rolls than respect from his oddball staff.",
	}

	err := tvShowService.Save(tvShow)
	assert.Nil(t, err)
}

func TestFailSave(t *testing.T) {
	tvShowRepo := mocks.NewInmemTVShowRepository()
	tvShowService := service.NewTVShowService(tvShowRepo)
	log.SetOutput(ioutil.Discard) // disable logrus for tests

	tvShow := &models.TVShow{}

	err := tvShowService.Save(tvShow)
	assert.NotNil(t, err)
}
