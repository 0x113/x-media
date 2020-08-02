package handler

import "github.com/0x113/x-media/tvshow/models"

// NOTE: these models are only for docs, they are not used in the handlers
type tvShowListResponse struct {
	TVShows []*models.TVShow `json:"tv_shows"`
}

type tvShowUpdateResponse struct {
	Errors       map[string]string `json:"errors"`
	UpdatedShows map[string]string `json:"updated_shows"`
}
