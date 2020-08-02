package handler

import (
	"fmt"
	"net/http"

	"github.com/0x113/x-media/movie-svc/models"
	"github.com/0x113/x-media/movie-svc/service"

	"github.com/go-openapi/runtime/middleware"
	"github.com/labstack/echo"
)

type movieHandler struct {
	movieService service.MovieService
}

// NewMovieHandler initiates the movie handlers
func NewMovieHandler(router *echo.Echo, movieService service.MovieService) {
	h := &movieHandler{movieService}
	// swagger
	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(opts, nil)
	router.File("/swagger.yaml", "./docs/swagger.yaml")
	router.GET("/docs", echo.WrapHandler(sh))

	router.POST("/api/v1/movies/update/all", h.UpdateAllMovies)
	router.GET("/api/v1/movies/all", h.GetAllMovies)
	router.GET("/api/v1/movies/:id", h.GetMovieByID)
}

// @Summary Update all movies
// @Description Calls the TMDb API to get data about movies from provided directories and saves it to the database
// @ID update-all-movies
// @Accept  json
// @Produce  json
// @Param name body updateAllMoviesPayload true "the language in which to update the movie data"
// @Success 200 {object} updateAllMoviesResponse
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /update/all [post]
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

// @Summary Get all movies
// @Description Retruns all movies from the database
// @ID get-all-movies
// @Produce  json
// @Success 200 {object} movieListResponse
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /get/all [get]
// GetAllMovies calls the service to get all movies from the database
func (h *movieHandler) GetAllMovies(c echo.Context) error {
	errMsg := new(models.Error)
	movies, err := h.movieService.GetAllMovies()
	if err != nil {
		errMsg.Code = http.StatusInternalServerError
		errMsg.Message = err.Error()
		c.JSON(errMsg.Code, errMsg)
		return err
	}

	res := map[string]interface{}{
		"movies": movies,
	}
	return c.JSON(http.StatusOK, res)
}

// GetMovieByID calls the movie service to get movie based on its id]
func (h *movieHandler) GetMovieByID(c echo.Context) error {
	errMsg := new(models.Error)
	id := c.Param("id")
	movie, err := h.movieService.GetMovieByID(id)
	if err != nil {
		errMsg.Code = http.StatusInternalServerError
		errMsg.Message = err.Error()
		c.JSON(errMsg.Code, errMsg)
		return err
	}
	fmt.Println(movie)

	return c.JSON(http.StatusOK, movie)
}
