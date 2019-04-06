package video

import (
	"encoding/json"
	"net/http"
)

type VideoHandler interface {
	// SaveMovies allows user to save movies from disk to the database
	UpdateMovies(w http.ResponseWriter, r *http.Request)
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
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Successfully saved!"})
}
