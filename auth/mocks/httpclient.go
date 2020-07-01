package mocks

import "net/http"

// MockClient defines mocked http client
type MockClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

// Do defines mocked Do func from the http client
func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}
