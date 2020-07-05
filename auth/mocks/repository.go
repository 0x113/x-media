package mocks

import (
	"fmt"

	"github.com/0x113/x-media/auth/models"
)

// MockAuthRepository represents in-memory authentication repository
type MockAuthRepository struct {
	tokens map[string]string
}

// NewMockAuthRepository creates new instance of the mocked auth repository
func NewMockAuthRepository() *MockAuthRepository {
	var tokens = map[string]string{}
	tokens["4fc2fa05-869d-45e3-aba8-09736dbb97d5"] = "Test"
	return &MockAuthRepository{tokens}
}

// Save the token in memory
func (m *MockAuthRepository) Save(username string, token *models.TokenDetails) error {
	if _, ok := m.tokens[token.AccessUuid]; ok {
		return fmt.Errorf("Unable to save the access token")
	}
	m.tokens[token.AccessUuid] = username

	if _, ok := m.tokens[token.RefreshUuid]; ok {
		return fmt.Errorf("Unable to save the refresh token")
	}
	m.tokens[token.RefreshUuid] = username

	return nil
}
