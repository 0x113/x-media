package auth

// User represents user model.
type User struct {
	ID       int64  `json:"-"`
	Username string `json:"username"`
	Password string `json:"password"`
}
