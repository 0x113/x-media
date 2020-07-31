package models

// swagger:response errorMsg
//
// Error defines the response error
type Error struct {
	Code    int    `json:"code" example:"500"`
	Message string `json:"message" example:"Couldn't get data from the TVmaze API"`
}
