package auth

// AuthRepository represtents authentication repository.
type AuthRepository interface {
	Create(user *User) error
	GetUser(username string) (*User, error)
}
