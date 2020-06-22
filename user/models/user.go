package models

// User information
type User struct {
	ID       int    `json:"id" validate:"omitempty"`
	Username string `json:"username" validate:"required,min=3"`
	Password string `json:"password" validate:"required,gte=8"`
	IsAdmin  bool   `json:"is_admin" validate:"isdefault"`
}
