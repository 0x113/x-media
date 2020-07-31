package handler

// tvShowNamePayload represents request body that should be sent
type tvShowNamePayload struct {
	Name string `json:"name" example:"BoJack Horseman"`
}
