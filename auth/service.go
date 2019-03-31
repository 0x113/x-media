package auth

// AuthService is the interface that provides authentication methods.
type AuthService interface {
	CreateUser(user *User) error
}

type authService struct {
	repo AuthRepository
}

// NewAuthService returns a new instance of authentication service.
func NewAuthService(repo AuthRepository) AuthService {
	return &authService{
		repo,
	}
}

func (s *authService) CreateUser(user *User) error {
	return s.repo.Create(user)
}
