package handler

import "github.com/0x113/x-media/movie-svc/models"

// NOTE: it's not used in the handler implementation, it's used for docs only

type updateAllMoviesPayload struct {
	Language string `json:"language" example:"en"`
}

type updateAllMoviesResponse struct {
	Errors       []*updateError  `json:"errors"`
	UpdateMovies []*updatedMovie `json:"updated_movies"`
}

type updateError struct {
	DirPath string `json:"/home/0x113/Movies/Heat.1995.mp4" example:"Unable to fund movie with such title"`
}

type updatedMovie struct {
	DirPath string `json:"/home/0x113/Movies/K-PAX.2001.mp4" example:"K-PAX"`
}

type movieListResponse struct {
	Movies []*models.Movie `json:"movies"`
}
