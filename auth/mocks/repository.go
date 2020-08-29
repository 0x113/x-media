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
	tokens["b66a7219-f07f-49cf-8163-189da2f5c8cc"] = "JohnDoe"
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

// Delete the token from memory
func (m *MockAuthRepository) Delete(uuid string) error {
	if _, ok := m.tokens[uuid]; !ok {
		return fmt.Errorf("There is no token with UUID: %s", uuid)
	}

	delete(m.tokens, uuid)
	return nil
}
