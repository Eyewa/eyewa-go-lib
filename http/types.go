package http

import "net/http"

type HTTPClient interface {
	Get(url, query string) (*http.Response, error)
	Post(url, query string, body interface{}) (*http.Response, error)
	NewRequest(method, url string, body interface{}) (*http.Request, error)
	Do(req *http.Request) (*http.Response, error)
	WithAuth() *client
	GetURL(path, query string) (string, error)
}

type client struct {
	http.Client
	baseUrl, authSecret string
	authRequired        bool
}
