package mocks

import "net/http"

// MockClient represents mocked http client
type MockClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

// Do is mocked function from the http.Client
func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	if m.DoFunc != nil {
		return m.DoFunc(req)
	}
	return &http.Response{}, nil
}
