package data

import "github.com/0x113/x-media/movie-svc/models"

// MovieRepository contains all methods for operation on the Movie model
type MovieRepository interface {
	Save(movie *models.Movie) error
	Update(movie *models.Movie) error
	GetByTitle(title string) (*models.Movie, error)
	GetByOriginalTitle(title string) (*models.Movie, error)
	GetAll() ([]*models.Movie, error)
}
