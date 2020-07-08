package service

import (
	"github.com/0x113/x-media/movie-svc/data"
	"github.com/0x113/x-media/movie-svc/external/tmdb"
	"github.com/0x113/x-media/movie-svc/httpclient"
	"github.com/0x113/x-media/movie-svc/models"

	log "github.com/sirupsen/logrus"
)

// MovieService defines the movie service
type MovieService interface {
	UpdateMovieByID(id int, lang, filePath string) error
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
func (s *movieService) UpdateMovieByID(id int, lang, filePath string) error {
	tmdbApiClient := &tmdb.TMDbAPIClient{s.httpClient}
	tmdbMovie, err := tmdbApiClient.GetTMDbMovieInfo(id, lang)
	if err != nil {
		return err
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

	// check if movie exits in the database NOTE: maybe getting by TMDb's ID is better ? ¯\_(ツ)_/¯
	// get movie based on it's title
	dbMovie, err := s.repo.GetByTitle(movie.Title) // NOTE: it's 11:29 PM CET and I have no idea how to handle this error
	if dbMovie == nil {
		if err := s.repo.Save(movie); err != nil {
			log.Errorf("Couldn't save new movie [%s]: %v", movie.Title, err)
			return err
		}
		log.Infof("Successfully saved new movie [%s]", movie.Title)
	} else {
		movie.ID = dbMovie.ID
		if err := s.repo.Update(movie); err != nil {
			log.Errorf("Couldn't update movie [%s]: %v", movie.Title, err)
			return err
		}
		log.Infof("Successfully updated movie [%s]", movie.Title)
	}
	return nil
}

// GetLocalTMDbID gets files from the given directory and calls the
// TMDb API to get ID of that movie.
func (s *movieService) GetLocalTMDbID(filename string) (int, error) {
	tmdbApiClient := &tmdb.TMDbAPIClient{s.httpClient}
	tmdbQMovie, err := tmdbApiClient.GetTMDbQueryMovieInfo(filename, "en") // NOTE: "lang" param is probably useless
	if err != nil {
		return 0, err
	}
	return tmdbQMovie.ID, nil
}
