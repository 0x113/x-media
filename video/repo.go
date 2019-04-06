package video

type VideoRepository interface {
	SaveMovie(movie *Movie) error
	FindAllMovies() ([]*Movie, error)
}
