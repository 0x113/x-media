package service

import (
	"fmt"
	"sync"

	"github.com/0x113/x-media/movie-svc/common"
	"github.com/0x113/x-media/movie-svc/data"
	"github.com/0x113/x-media/movie-svc/external/tmdb"
	"github.com/0x113/x-media/movie-svc/httpclient"
	"github.com/0x113/x-media/movie-svc/models"
	"github.com/0x113/x-media/movie-svc/utils/filenameparser"
	"github.com/0x113/x-media/movie-svc/utils/scandir"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MovieService defines the movie service
type MovieService interface {
	UpdateMovieByID(id int, lang, filePath string, mutex *sync.Mutex) (*models.Movie, error)
	UpdateAllMovies(lang string) (map[string]string, map[string]string)
	GetAllMovies() ([]*models.Movie, error)
	GetLocalTMDbID(filename string) (int, error)
	GetMovieByID(id string) (*models.Movie, error)
}

type movieService struct {
	repo       data.MovieRepository
	httpClient httpclient.HTTPClient
}

// NewMovieService returns new insance of the movie service
func NewMovieService(repo data.MovieRepository, httpClient httpclient.HTTPClient) MovieService {
	return &movieService{repo, httpClient}
}

// UpdateMovieByID calls the TMDb API to get data about movie
// based on its ID and saves it to the database if doesn't exist
// or updates if exists.
func (s *movieService) UpdateMovieByID(id int, lang, filePath string, mutex *sync.Mutex) (*models.Movie, error) {
	tmdbAPIClient := &tmdb.TMDbAPIClient{s.httpClient}
	tmdbMovie, err := tmdbAPIClient.GetTMDbMovieInfo(id, lang)
	if err != nil {
		log.Errorf("Unable to get the data from the TMDb API [movie_id: %d, lang: %s]: %v", id, lang, err)
		return nil, err
	}

	var genres []string
	for _, g := range tmdbMovie.Genres {
		genres = append(genres, g.Name)
	}

	movie := &models.Movie{
		TMDbID:           tmdbMovie.ID,
		IMDbID:           tmdbMovie.IMDbID,
		Title:            tmdbMovie.Title,
		Overview:         tmdbMovie.Overview,
		OriginalTitle:    tmdbMovie.OriginalTitle,
		OriginalLanguage: tmdbMovie.OriginalLanguage,
		ReleaseDate:      tmdbMovie.ReleaseDate,
		Genres:           genres,
		Rating:           tmdbMovie.VoteAverage,
		VoteCount:        tmdbMovie.VoteCount,
		Runtime:          tmdbMovie.Runtime,
		BackdropPath:     tmdbMovie.BackdropPath,
		PosterPath:       tmdbMovie.PosterPath,
		DirPath:          filePath,
	}

	mutex.Lock()
	// check if movie exists in the database NOTE: maybe getting by TMDb's ID is better ? ¯\_(ツ)_/¯
	// get movie based on it's title
	dbMovie, err := s.repo.GetByOriginalTitle(movie.OriginalTitle) // NOTE: it's 11:29 PM CET and I have no idea how to handle this error
	if dbMovie == nil {
		movie.ID = primitive.NewObjectID()
		if err := s.repo.Save(movie); err != nil {
			log.Errorf("Couldn't save new movie [%s]: %v", movie.Title, err)
			return nil, err
		}
		log.Infof("Successfully saved new movie [%s]", movie.Title)
	} else {
		movie.ID = dbMovie.ID
		if err := s.repo.Update(movie); err != nil {
			log.Errorf("Couldn't update movie [%s]: %v", movie.Title, err)
			return nil, err
		}
		log.Infof("Successfully updated movie [%s]", movie.Title)
	}

	mutex.Unlock()
	return movie, nil
}

// UpdateAllMovies scans the given directories for video files like mp4, mkv etc.
// Then it calls the TMDb API to get data about every single one and saves new movie
// to the database if it doesn't exist or updates movie if there is already one.
func (s *movieService) UpdateAllMovies(lang string) (map[string]string, map[string]string) {
	errorsMap := make(map[string]string)
	type moviePathID struct {
		filepath string
		id       int
	}
	var movieIDs []*moviePathID // contains list of moviePathID (filepath: tmdb_id)

	for _, dir := range common.Config.MovieDirectories {
		// get files from the given directories
		files, err := scandir.GetFiles(dir, []string{".mp4", ".mkv"})
		if err != nil {
			errorsMap[dir] = err.Error()
		}

		// for every single file parse filename to get movie title and
		// send request to the TMDb API to get movie id
		// FIXME: error handling like 401 from TMDb's API
		for _, f := range files {
			title, err := filenameparser.CreateTitle(f)
			if err != nil {
				errorsMap[f] = err.Error()
			}
			id, err := s.GetLocalTMDbID(title)
			if err != nil {
				errorsMap[title] = err.Error()
				continue
			}
			movieIDs = append(movieIDs, &moviePathID{f, id})
		}
	}

	updatedMovies := make(map[string]string)
	var wg sync.WaitGroup
	var mutex sync.Mutex
	wg.Add(len(movieIDs))

	for _, m := range movieIDs {
		go func(m *moviePathID) {
			defer wg.Done()
			movie, err := s.UpdateMovieByID(m.id, lang, m.filepath, &mutex)
			if err != nil {
				mutex.Lock()
				errorsMap[m.filepath] = err.Error()
				mutex.Unlock()
				return
			}
			mutex.Lock()
			updatedMovies[movie.DirPath] = movie.Title
			mutex.Unlock()
		}(m)
	}
	wg.Wait()
	return updatedMovies, errorsMap
}

// GetAllMovies calls the database layer and returns all movies from the database
func (s *movieService) GetAllMovies() ([]*models.Movie, error) {
	movies, err := s.repo.GetAll()
	if err != nil {
		log.Errorf("Couldn't get movies from the database: %v", err)
		return nil, fmt.Errorf("Couldn't get movies from the database")
	}

	log.Infoln("Successfully found all movies")
	return movies, nil
}

// GetLocalTMDbID calls the TMDb API to get movie ID based on its title.
func (s *movieService) GetLocalTMDbID(title string) (int, error) {
	tmdbAPIClient := &tmdb.TMDbAPIClient{s.httpClient}
	tmdbQMovie, err := tmdbAPIClient.GetTMDbQueryMovieInfo(title, "en") // NOTE: "lang" param is probably useless
	if err != nil {
		return 0, err
	}
	return tmdbQMovie.ID, nil
}

// GetMovieByID converts provided id to the ObjectID and then
// returns movie from the database with this id
func (s *movieService) GetMovieByID(id string) (*models.Movie, error) {
	movieID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Errorf("Unable to convert string id to ObjectID: %s", id)
		return nil, fmt.Errorf("Couldn't get movie from the database: unable to convert provided id")
	}

	movie, err := s.repo.GetByID(movieID)
	if err != nil {
		log.Errorf("Unable to get movie by id: %s; err: %v", id, err)
		return nil, fmt.Errorf("Couldn't get movie from the database")
	}

	log.Infof("Successfully found movie with id: %s", id)
	return movie, nil
}
