package models

import "time"

// User information
type User struct {
	ID        int       `json:"id" validate:"omitempty"`
	Username  string    `json:"username" validate:"required,min=3,max=32"`
	Password  string    `json:"password" validate:"required,gte=8"`
	IsAdmin   bool      `json:"is_admin" validate:"isdefault"`
	CreatedAt time.Time `json:"created_at" validate:"isdefault"`
	UpdatedAt time.Time `json:"updated_at" validate:"isdefault"`
}
