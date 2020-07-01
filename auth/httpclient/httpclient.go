package httpclient

import "net/http"

// HTTPClient defines the http client
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}
