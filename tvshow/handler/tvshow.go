package handler

import (
	"net/http"

	"github.com/0x113/x-media/tvshow/models"
	"github.com/0x113/x-media/tvshow/service"

	"github.com/go-openapi/runtime/middleware"
	"github.com/labstack/echo"
)

type tvShowHandler struct {
	tvShowService service.TVShowService
}

// NewTVShowHandler initiates tv show handlers
func NewTVShowHandler(router *echo.Echo, tvShowService service.TVShowService) {
	handler := &tvShowHandler{tvShowService}

	// swagger
	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(opts, nil)
	router.File("/swagger.yaml", "./docs/swagger.yaml")
	router.GET("/docs", echo.WrapHandler(sh))

	// TODO: use groups
	router.POST("/api/v1/tvshows/get", handler.GetTVShow)
	router.GET("/api/v1/tvshows/get/all", handler.GetAllTVShows)
	router.GET("/api/v1/tvshows/update/all", handler.UpdateAllTVShows)
}

// @Summary Get tv show
// @Description Returns tv shows
// @ID get-tvshow-by-name
// @Accept  json
// @Produce  json
// @Param name body tvShowNamePayload true "title of the tv show"
// @Success 200 {object} models.TVShow
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /get [post]
// GetTVShow calls service layer to get an existing tv show from the database
func (h *tvShowHandler) GetTVShow(c echo.Context) error {
	errMsg := &models.Error{} // error message
	payload := new(tvShowNamePayload)
	if err := c.Bind(&payload); err != nil {
		errMsg.Code = http.StatusBadRequest
		errMsg.Message = err.Error()
		c.JSON(errMsg.Code, errMsg)
		return err
	}

	tvShow, err := h.tvShowService.GetTVShowByName(payload.Name)
	if err != nil {
		errMsg.Code = http.StatusNotFound
		errMsg.Message = err.Error()
		c.JSON(errMsg.Code, errMsg)
		return err
	}
	return c.JSON(http.StatusOK, tvShow)
}

// @Summary Update all tv shows
// @Description Calls the third party API (TVMaze at this moment) to get data about tv shows from the local drive
// @ID update-all-tv-shows
// @Produce  json
// @Success 200 {object} tvShowUpdateResponse
// @Router /update/all [get]
// UpdateAllTVShows calls service layer and updates all tv shows which
// are in specified dirs
func (h *tvShowHandler) UpdateAllTVShows(c echo.Context) error {
	response := make(map[string]interface{})
	updatedShows, errList := h.tvShowService.UpdateAllTVShows()

	response["errors"] = errList
	response["updated_shows"] = updatedShows
	return c.JSON(http.StatusOK, response)
}

// @Summary Get all tv shows
// @Description Returns all the tv shows from the database
// @ID get-all-tv-shows
// @Produce json
// @Success 200 {object} tvShowListResponse
// @Failure 500 {object} models.Error
// @Router /get/all [get]
// GetAllTVShows calls service layer and returns all tv shows from the database
func (h *tvShowHandler) GetAllTVShows(c echo.Context) error {
	errMsg := &models.Error{}
	tvShows, err := h.tvShowService.GetAllTVShows()
	if err != nil {
		errMsg.Code = http.StatusInternalServerError
		errMsg.Message = err.Error()
		c.JSON(errMsg.Code, errMsg)
		return err
	}

	msg := map[string]interface{}{
		"tv_shows": tvShows,
	}
	return c.JSON(http.StatusOK, msg)
}
