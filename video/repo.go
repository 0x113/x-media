package video

type VideoRepository interface {
	// SaveMovie saves movie to the database
	SaveMovie(movie *Movie) error
	// RemoveMovieByFileName removes movie from the database
	RemoveMovieByFileName(movieFileName string) error
	// FindAllMovies returns list of all movies from the database
	FindAllMovies() ([]*Movie, error)
	// GetMovieById returns movie with certain id
	GetMovieById(id string) (*Movie, error)
	// SaveTvSeries saves tv series to the database
	SaveTvSeries(tvSeries *TVSeries) error
	// RemoveTvSeriesByDirName removes tv seris from the database
	RemoveTvSeriesByDirName(dirName string) error
	// FindAllTvSeries returns list of all tv series from the database
	FindAllTvSeries() ([]*TVSeries, error)
	// GetTvSeriesById returns tv series with certain id
	GetTvSeriesById(id string) (*TVSeries, error)
}
