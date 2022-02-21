package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/eyewa/eyewa-go-lib/utils"
)

const DefaultContentType = "application/json"

func (c *client) WithAuth() *client {
	c.authRequired = true

	return c
}

func (c *client) GetURL(path, query string) (string, error) {
	reqUrl, err := url.Parse(c.baseUrl)
	if err != nil {
		return "", err
	}

	reqUrl.Path = path
	if query != "" {
		reqUrl.RawQuery = query
	}

	return reqUrl.String(), nil
}

func (c *client) Get(path, query string) (*http.Response, error) {
	reqUrl, err := c.GetURL(path, query)
	if err != nil {
		return nil, err
	}

	req, err := c.NewRequest("GET", reqUrl, nil)
	if err != nil {
		return nil, err
	}

	return c.Do(req)
}

func (c *client) Post(path, query string, body interface{}) (*http.Response, error) {
	reqUrl, err := c.GetURL(path, query)
	if err != nil {
		return nil, err
	}

	req, err := c.NewRequest("POST", reqUrl, body)
	if err != nil {
		return nil, err
	}

	return c.Do(req)
}

func (c *client) Do(req *http.Request) (*http.Response, error) {
	return c.Client.Do(req)
}

func (c *client) NewRequest(method, url string, body interface{}) (*http.Request, error) {
	var reqBody io.Reader

	if body != nil {
		bodyMarshalled, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}

		reqBody = bytes.NewBuffer(bodyMarshalled)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", DefaultContentType)

	if c.authRequired {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.authSecret))
	}

	return req, nil
}

func NewClient(baseUrl, authSecret string) HTTPClient {
	return &client{
		baseUrl:    baseUrl,
		authSecret: authSecret,
		Client:     *utils.GetHTTPClient(),
	}
}
