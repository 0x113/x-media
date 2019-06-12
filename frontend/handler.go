package frontend

import (
	"encoding/json"
	"net/http"
)

type FrontendHandler interface {
	ServeFrontend(w http.ResponseWriter, r *http.Request)
}

type frontendHandler struct {
	frontendService FrontendService
}

func NewFrontendHandler(frontendService FrontendService) FrontendHandler {
	return &frontendHandler{
		frontendService,
	}
}

func (h *frontendHandler) ServeFrontend(w http.ResponseWriter, r *http.Request) {
	dir, err := h.frontendService.FrontendDir()
	if err != nil {
		response := make(map[string]interface{})
		response["error"] = "Frontend dir doesn't exist"
		json.NewEncoder(w).Encode(response)
		return
	}
	http.ServeFile(w, r, dir)
	//	http.FileServer(http.Dir(dir))
}
