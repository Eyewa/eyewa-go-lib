package http

import "net/http"

type HTTPClient interface {
	Get(url string) (*http.Response, error)
	Post(url string, body interface{}) (*http.Response, error)
	NewRequest(method, url string, body interface{}) (*http.Request, error)
	Do(req *http.Request) (*http.Response, error)
	WithAuth() *client
	GetURLWithPath(path string) (string, error)
}

type client struct {
	http.Client
	baseUrl, authSecret string
	authRequired        bool
}
