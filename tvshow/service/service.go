package service

import (
	"fmt"
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
	UpdateAllTVShows() (map[string]string, map[string]string)
	UpdateTVShow(name string, mutex *sync.Mutex) (*models.TVShow, error)
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
func (s *tvShowService) UpdateTVShow(dirPath string, mutex *sync.Mutex) (*models.TVShow, error) {
	nameSplit := strings.Split(dirPath, "/")
	name := createName(nameSplit[len(nameSplit)-1])
	// get tv show data from TVmaze API
	tvMazeInfo, err := tvmaze.GetTVmazeTVShowInfo(s.client, name)
	if err != nil {
		return nil, err
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
		DirPath:   dirPath,
	}
	// validate new TVShow object
	validate := validator.New()
	if err := validate.Struct(tvShow); err != nil {
		log.Errorf("Couldn't validate tv show[dir=%s]; err: %v", name, err)
		return nil, fmt.Errorf("Couldn't validate tv show [dir=%s]", name)
	}

	mutex.Lock()
	existingShow, err := s.tvShowRepo.GetByName(tvShow.Name)

	if existingShow == nil {
		if err := s.Save(tvShow); err != nil { // NOTE: here tv show is validated twice, need to be changed
			log.Debugf("Couldn't save new tv show[%s]; err: %v", tvShow.Name, err)
			return nil, err
		}
		log.Infof("Successfully saved new tv show[%s]", tvShow.Name)
	} else {
		tvShow.ID = existingShow.ID
		if err := s.tvShowRepo.Update(tvShow); err != nil {
			log.Debugf("Couldn't update tv show[%s]; err: %v", tvShow.Name, err)
			return nil, err
		}
		log.Infof("Successfully updated tv show[%s]", tvShow.Name)
	}
	// unlock mutex
	mutex.Unlock()
	return tvShow, nil
}

// UpdateAllTVShows reads directory names, removes special char like "_,/"
// and calls TVmaze api to get data
func (s *tvShowService) UpdateAllTVShows() (map[string]string, map[string]string) {
	tvShowDirs := getDirectories()

	// get data from TVmaze api
	var wg sync.WaitGroup
	var mutex sync.Mutex
	wg.Add(len(tvShowDirs))

	updatedShows := make(map[string]string) // this var contains names of updated tv shows (for http handler)
	errorsMap := make(map[string]string)

	for i, n := range tvShowDirs {
		go func(i int, n string) {
			defer wg.Done()
			tvShow, err := s.UpdateTVShow(n, &mutex)
			if err != nil {
				mutex.Lock()
				errorsMap[n] = err.Error()
				mutex.Unlock()
				return
			}
			mutex.Lock() // lock mutex to avoid race condition
			updatedShows[tvShow.DirPath] = tvShow.Name
			mutex.Unlock()
		}(i, n)
	}
	wg.Wait()

	return updatedShows, errorsMap
}

// GetTVShowByName returns tv show if exists
func (s *tvShowService) GetTVShowByName(name string) (*models.TVShow, error) {
	tvShow, err := s.tvShowRepo.GetByName(name)
	if err != nil {
		log.Debugf("Unable to get tv show from the database; err: %v", err)
		return nil, fmt.Errorf("There is no %s in the tv show database", name)
	}
	log.Infof("Successfully found tv show: %s", name)
	return tvShow, nil
}

// GetAllTVShows returns all tv show from the database
func (s *tvShowService) GetAllTVShows() ([]*models.TVShow, error) {
	tvShows, err := s.tvShowRepo.GetAll()
	if err != nil {
		log.Debugf("Couldn't get all tv shows from the database; err: %v", err)
		return nil, fmt.Errorf("Couldn't get all tv shows from the database")
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
				if !strings.HasSuffix(dir, "/") {
					dir += "/"
				}
				tvShowDirs = append(tvShowDirs, dir+f.Name())
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
