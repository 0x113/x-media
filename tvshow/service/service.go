package service

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/0x113/x-media/tvshow/common"
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

// NewTVShowService creates new instance of TVShowService
func NewTVShowService(tvShowRepo data.TVShowRepository) TVShowService {
	return &tvShowService{tvShowRepo}
}

// Save calls the db layer to save tv show
func (s *tvShowService) Save(tvShow *models.TVShow) error {
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

// directoryExists checks if a directory exists and
// is not a file
func directoryExists(dirName string) bool {
	info, err := os.Stat(dirName)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

// getDirectories reads directory paths from config
// and returns list of directories inside those path
func getDirectories() []string {
	tvShowDirs := []string{} // contains folders name e.g. "The_Office"
	dirs := common.Config.TVShowDirectories

	for _, dir := range dirs {
		// check if directory exists
		if !directoryExists(dir) {
			log.Debugf("Directory [%s] does not exist or it is not a directory", dir)
			continue
		}

		files, err := ioutil.ReadDir(dir)
		if err != nil {
			log.Debugf("Couldn't read dir [%s]; err: %v", dir, err)
			continue
		}

		for _, f := range files {
			if f.IsDir() {
				tvShowDirs = append(tvShowDirs, f.Name())
			}
		}
	}

	return tvShowDirs
}

// createName removes special characters like "_, /" from directory name
// e.g. "The_Office" -> "The Office"
func createName(dirName string) string {
	charsToRemove := []string{".", ",", "_"}
	name := dirName
	for _, c := range charsToRemove {
		if strings.Contains(dirName, c) {
			name = strings.Replace(dirName, c, " ", -1)
		}
	}
	// remove "/" from end of the directory
	if strings.HasSuffix(name, "/") {
		name = strings.TrimSuffix(name, "/")
	}

	return name
}
