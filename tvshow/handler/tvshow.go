package handler

import (
	"net/http"

	"github.com/0x113/x-media/tvshow/service"

	"github.com/labstack/echo"
)

type tvShowHandler struct {
	tvShowService service.TVShowService
}

// NewTVShowHandler initiates tv show handlers
func NewTVShowHandler(router *echo.Echo, tvShowService service.TVShowService) {
	handler := &tvShowHandler{tvShowService}
	// TODO: use groups
	router.POST("/api/v1/tvshows/get", handler.GetTVShow)
	router.GET("/api/v1/tvshows/get/all", handler.GetAllTVShows)
	router.GET("/api/v1/tvshows/update/all", handler.UpdateAllTVShows)
}

// GetTVShow calls service layer to get an existing tv show from the database
func (h *tvShowHandler) GetTVShow(c echo.Context) error {
	var reqBody struct {
		Name string `json:"name"`
	}
	if err := c.Bind(&reqBody); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "error")
		return err
	}

	tvShow, err := h.tvShowService.GetTVShowByName(reqBody.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "error")
		return err
	}
	return c.JSON(http.StatusOK, tvShow)
}

// UpdateAllTVShows calls service layer and updates all tv shows which
// are in specified dirs
func (h *tvShowHandler) UpdateAllTVShows(c echo.Context) error {
	response := make(map[string]interface{})
	updatedShows, errList := h.tvShowService.UpdateAllTVShows()

	response["errors"] = errList
	response["updated_shows"] = updatedShows
	return c.JSON(http.StatusOK, response)
}

// GetAllTVShows calls service layer and returns all tv shows from the database
func (h *tvShowHandler) GetAllTVShows(c echo.Context) error {
	tvShows, err := h.tvShowService.GetAllTVShows()
	if err != nil {
		c.JSON(http.StatusInternalServerError, "error")
		return err
	}

	msg := map[string]interface{}{
		"tv_shows": tvShows,
	}
	return c.JSON(http.StatusOK, msg)
}
