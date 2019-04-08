package video

import (
	"encoding/json"
	"net/http"
)

type VideoHandler interface {
	// UpdateMovies allows user to save movies from disk to the database
	UpdateMovies(w http.ResponseWriter, r *http.Request)
	// AllMovies returns all movies in json format
	AllMovies(w http.ResponseWriter, r *http.Request)
	// UpdateTvSeries alloew user to save tv series from disk to the database
	UpdateTvSeries(w http.ResponseWriter, r *http.Request)
	// AllTvSeries returns all tv series in json format
	AllTvSeries(w http.ResponseWriter, r *http.Request)
}

type videoHandler struct {
	videoService VideoService
}

func NewVideoHandler(videoService VideoService) VideoHandler {
	return &videoHandler{
		videoService,
	}
}

func (h *videoHandler) UpdateMovies(w http.ResponseWriter, r *http.Request) {
	err := h.videoService.Save()
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Successfully saved"})

}

func (h *videoHandler) AllMovies(w http.ResponseWriter, r *http.Request) {
	movies, err := h.videoService.AllMovies()
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	json.NewEncoder(w).Encode(movies)
}

func (h *videoHandler) UpdateTvSeries(w http.ResponseWriter, r *http.Request) {
	err := h.videoService.SaveTVShows()
	response := make(map[string]string)
	if err != nil {
		response["error"] = err.Error()
		json.NewEncoder(w).Encode(response)
		return
	}
	response["message"] = "Successfully updated tv series"
	json.NewEncoder(w).Encode(response)
}

func (h *videoHandler) AllTvSeries(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})
	tvSeries, err := h.videoService.AllTvSeries()
	if err != nil {
		response["error"] = "Unable to get all tv series"
		json.NewEncoder(w).Encode(response)
		return
	}
	response["tv_series"] = tvSeries
	json.NewEncoder(w).Encode(response)
}
