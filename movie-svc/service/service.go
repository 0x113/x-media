package service

import (
	"sync"

	"github.com/0x113/x-media/movie-svc/common"
	"github.com/0x113/x-media/movie-svc/data"
	"github.com/0x113/x-media/movie-svc/external/tmdb"
	"github.com/0x113/x-media/movie-svc/httpclient"
	"github.com/0x113/x-media/movie-svc/models"
	"github.com/0x113/x-media/movie-svc/utils/filenameparser"
	"github.com/0x113/x-media/movie-svc/utils/scandir"

	log "github.com/sirupsen/logrus"
)

// MovieService defines the movie service
type MovieService interface {
	UpdateMovieByID(id int, lang, filePath string, mutex *sync.Mutex) (*models.Movie, error)
	UpdateAllMovies(lang string) (map[string]string, map[string]string)
	GetLocalTMDbID(filename string) (int, error)
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
	tmdbApiClient := &tmdb.TMDbAPIClient{s.httpClient}
	tmdbMovie, err := tmdbApiClient.GetTMDbMovieInfo(id, lang)
	if err != nil {
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
	dbMovie, err := s.repo.GetByTitle(movie.Title) // NOTE: it's 11:29 PM CET and I have no idea how to handle this error
	if dbMovie == nil {
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
	type moviePathId struct {
		filepath string
		id       int
	}
	var movieIDs []*moviePathId // contains list of moviePathId (filepath: tmdb_id)

	for _, dir := range common.Config.MovieDirectories {
		// get files from the given directories
		files, err := scandir.GetFiles(dir, []string{".mp4", ".mkv"})
		if err != nil {
			errorsMap[dir] = err.Error()
		}

		// for every single file parse filename to get movie title and
		// send request to the TMDb API to get movie id
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
			movieIDs = append(movieIDs, &moviePathId{f, id})
		}
	}

	updatedMovies := make(map[string]string)
	var wg sync.WaitGroup
	var mutex sync.Mutex
	wg.Add(len(movieIDs))

	for _, m := range movieIDs {
		go func(m *moviePathId) {
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
	return nil, nil
}

// GetLocalTMDbID calls the TMDb API to get movie ID based on its title.
func (s *movieService) GetLocalTMDbID(title string) (int, error) {
	tmdbApiClient := &tmdb.TMDbAPIClient{s.httpClient}
	tmdbQMovie, err := tmdbApiClient.GetTMDbQueryMovieInfo(title, "en") // NOTE: "lang" param is probably useless
	if err != nil {
		return 0, err
	}
	return tmdbQMovie.ID, nil
}
