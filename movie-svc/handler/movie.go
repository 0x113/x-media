package handler

import (
	"net/http"

	"github.com/0x113/x-media/movie-svc/models"
	"github.com/0x113/x-media/movie-svc/service"

	"github.com/labstack/echo"
)

type movieHandler struct {
	movieService service.MovieService
}

// NewMovieHandler initiates the movie handlers
func NewMovieHandler(router *echo.Echo, movieService service.MovieService) {
	h := &movieHandler{movieService}
	router.POST("/api/v1/movies/update/all", h.UpdateAllMovies)
}

// UpdateAllMovies calls the service to update all movies from the given directories
func (h *movieHandler) UpdateAllMovies(c echo.Context) error {
	var reqBody struct {
		Language string `json:"language"`
	}
	if err := c.Bind(&reqBody); err != nil {
		errMsg := models.Error{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		c.JSON(errMsg.Code, errMsg)
		return err
	}

	response := make(map[string]interface{})
	updatedMovies, errList := h.movieService.UpdateAllMovies(reqBody.Language)

	response["errors"] = errList
	response["updated_movies"] = updatedMovies

	return c.JSON(http.StatusOK, response)
}
