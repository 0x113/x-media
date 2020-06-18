package service

import (
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"github.com/0x113/x-media/tvshow/common"
	"github.com/0x113/x-media/tvshow/data"
	"github.com/0x113/x-media/tvshow/external/tvmaze"
	"github.com/0x113/x-media/tvshow/models"
	"github.com/0x113/x-media/tvshow/utils"

	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TVShowService describes a tv show service
type TVShowService interface {
	Save(tvShow *models.TVShow) error
	UpdateTVShows() error
	UpdateTVShow(name string, mutex *sync.Mutex, wg *sync.WaitGroup) error
	GetTVShowByName(name string) (*models.TVShow, error)
	GetAllTVShows() ([]*models.TVShow, error)
}

type tvShowService struct {
	client     utils.HttpClient
	tvShowRepo data.TVShowRepository
}

// NewTVShowService creates new instance of TVShowService
func NewTVShowService(client utils.HttpClient, tvShowRepo data.TVShowRepository) TVShowService {
	return &tvShowService{client, tvShowRepo}
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

// UpdateTVShow reads directory names, removes special char like "_,/"
// and calls tvmaze api to get data
func (s *tvShowService) UpdateTVShow(name string, mutex *sync.Mutex, wg *sync.WaitGroup) error {
	// get tv show data from TVmaze API
	tvMazeInfo, err := tvmaze.GetTVmazeTVShowInfo(s.client, name)
	if err != nil {
		return err
	}
	// create new TVShow object
	tvShow := &models.TVShow{
		Name:      tvMazeInfo.Show.Name,
		Language:  tvMazeInfo.Show.Language,
		Genres:    tvMazeInfo.Show.Genres,
		Runtime:   tvMazeInfo.Show.Runtime,
		Premiered: tvMazeInfo.Show.Premiered,
		Rating:    tvMazeInfo.Show.Rating.Average,
		PosterURL: tvMazeInfo.Show.Image.Original,
		Summary:   tvMazeInfo.Show.Summary,
	}
	// validate new TVShow object
	validate := validator.New()
	if err := validate.Struct(tvShow); err != nil {
		log.Errorf("Couldn't validate tv show; err: %v", err)
		return err
	}

	mutex.Lock()
	existingShow, err := s.tvShowRepo.GetByName(tvShow.Name)

	if existingShow == nil {
		if err := s.Save(tvShow); err != nil { // NOTE: here tv show is validated twice, need to be changed
			log.Debugf("Couldn't save new tv show[%s]; err: %v", tvShow.Name, err)
			return err
		}
		log.Infof("Successfully saved new tv show[%s]", tvShow.Name)
	} else {
		tvShow.ID = existingShow.ID
		if err := s.tvShowRepo.Update(tvShow); err != nil {
			log.Debugf("Couldn't update tv show[%s]; err: %v", tvShow.Name, err)
			return err
		}
		log.Infof("Successfully updated tv show[%s]", tvShow.Name)
	}
	// unlock mutex and finish wait group
	mutex.Unlock()
	wg.Done()
	return nil
}

// UpdateTVShows reads directory names, removes special char like "_,/"
// and calls tvmaze api to get data
func (s *tvShowService) UpdateTVShows() error {
	// get tv show names
	tvShowDirs := getDirectories()
	names := []string{}
	for _, dir := range tvShowDirs {
		names = append(names, createName(dir))
	}

	// get data from TVmaze api
	var wg sync.WaitGroup
	var mutex sync.Mutex
	wg.Add(len(names))

	finished := make(chan bool, 1)
	errs := make(chan error)
	for _, n := range names {
		go func(n string) {
			if err := s.UpdateTVShow(n, &mutex, &wg); err != nil {
				errs <- err
			}
		}(n)
	}

	go func() {
		wg.Wait()
		close(finished)
	}()

	select {
	case <-finished:
	case err := <-errs:
		close(errs)
		return err
	}
	return nil
}

// GetTVShowByName returns tv show if exists
func (s *tvShowService) GetTVShowByName(name string) (*models.TVShow, error) {
	tvShow, err := s.tvShowRepo.GetByName(name)
	if err != nil {
		log.Debugf("Unable to get tv show from the database; err: %v", err)
		return nil, err
	}
	log.Infof("Successfully found tv show: %s", name)
	return tvShow, nil
}

// GetAllTVShows returns all tv show from the database
func (s *tvShowService) GetAllTVShows() ([]*models.TVShow, error) {
	tvShows, err := s.tvShowRepo.GetAll()
	if err != nil {
		log.Debugf("Couldn't get all tv shows from the database; err: %v", err)
		return nil, err
	}

	log.Infof("Successfully found all the shows")
	return tvShows, nil
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
