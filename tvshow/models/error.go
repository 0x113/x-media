package models

// Error defines the response error
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
