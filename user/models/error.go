package models

// Error defines the response error
type Error struct {
	Code    int    `json:"code" example:"500"`
	Message string `json:"message" example:"Server error"`
}
