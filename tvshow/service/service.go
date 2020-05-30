package service

import (
	"github.com/0x113/x-media/tvshow/data"
	"github.com/0x113/x-media/tvshow/models"

	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TVShowService describes a tv show service
type TVShowService interface {
	Save(tvShow *models.TVShow) error
}

type tvShowService struct {
	tvShowRepo data.TVShowRepository
}

func NewTVShowService(tvShowRepo data.TVShowRepository) TVShowService {
	return &tvShowService{tvShowRepo}
}

// Save calls the db layer to save tv show
func (s *tvShowService) Save(tvShow *models.TVShow) error {
	// TODO: validate
	validate := validator.New()
	if err := validate.Struct(tvShow); err != nil {
		log.Errorf("Couldn't validate tv show; err: %v", err)
		return err
	}
	// create new object id
	tvShow.ID = primitive.NewObjectID()

	if err := s.tvShowRepo.Save(tvShow); err != nil {
		log.Errorf("Couldn't save tv show [name=%s]; err: %v", tvShow.Name, err)
		return err
	}
	return nil
}
