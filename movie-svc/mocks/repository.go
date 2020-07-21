package mocks

import (
	"fmt"

	"github.com/0x113/x-media/movie-svc/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockMovieRepository represents in-memory user repository
type MockMovieRepository struct {
	movies map[string]*models.Movie
}

// NewMockMovieRepository creates new mocked movie repository
func NewMockMovieRepository() *MockMovieRepository {
	var movies = map[string]*models.Movie{}
	id, err := primitive.ObjectIDFromHex("507f1f77bcf86cd799439011")
	if err != nil {
		panic(err)
	}
	movies["Heat"] = &models.Movie{
		ID:               id,
		TMDbID:           949,
		Title:            "Heat",
		Overview:         "Obsessive master thief, Neil McCauley leads a top-notch crew on various daring heists throughout Los Angeles while determined detective, Vincent Hanna pursues him without rest. Each man recognizes and respects the ability and the dedication of the other even though they are aware their cat-and-mouse game may end in violence.",
		OriginalTitle:    "Heat",
		OriginalLanguage: "en",
		ReleaseDate:      "1995-12-15",
		Genres: []string{
			"Action",
			"Crime",
			"Drama",
			"Thriller",
		},
		Rating:       7.9,
		Runtime:      170,
		BackdropPath: "/rfEXNlql4CafRmtgp2VFQrBC4sh.jpg",
		PosterPath:   "/rrBuGu0Pjq7Y2BWSI6teGfZzviY.jpg",
		DirPath:      "/home/y0x/Videos/Heat.1995.mp4",
	}
	return &MockMovieRepository{movies}
}

// Save movie in memory
func (m *MockMovieRepository) Save(movie *models.Movie) error {
	if _, ok := m.movies[movie.Title]; ok {
		return fmt.Errorf("Couldn't save movie %s: it already exists in the database", movie.Title)
	}

	m.movies[movie.Title] = movie
	return nil
}

// Update movie in memory
func (m *MockMovieRepository) Update(movie *models.Movie) error {
	if _, ok := m.movies[movie.Title]; !ok {
		return fmt.Errorf("Couldn't update movie %s: no such movie in the database", movie.Title)
	}

	m.movies[movie.Title] = movie
	return nil
}

// GetByTitle returns movie from the mocked database if it exists
func (m *MockMovieRepository) GetByTitle(title string) (*models.Movie, error) {
	if _, ok := m.movies[title]; ok {
		return m.movies[title], nil
	}

	return nil, fmt.Errorf("Couldn't get movie %s: no such movie in the database", title)
}

// GetByOriginalTitle returns movie from the mocked database based on its original title
func (m *MockMovieRepository) GetByOriginalTitle(title string) (*models.Movie, error) {
	for _, movie := range m.movies {
		if movie.OriginalTitle == title {
			return movie, nil
		}
	}

	return nil, fmt.Errorf("Couldn't get movie %s: no such movie in the database", title)
}

// GetAll returns all movies from the mocked database
func (m *MockMovieRepository) GetAll() ([]*models.Movie, error) {
	var movies []*models.Movie
	for _, movie := range m.movies {
		movies = append(movies, movie)
	}

	return movies, nil
}

// GetByID returns movie from the mocked database by its id
func (m *MockMovieRepository) GetByID(id primitive.ObjectID) (*models.Movie, error) {
	for _, movie := range m.movies {
		if movie.ID == id {
			return movie, nil
		}
	}
	return nil, fmt.Errorf("Unable to find movie with id: %s", id)
}
