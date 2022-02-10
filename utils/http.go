package utils

import (
	"net/http"
	"time"
)

func GetHTTPClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives:     true,
			MaxIdleConns:          1,
			MaxIdleConnsPerHost:   1,
			TLSHandshakeTimeout:   30 * time.Second,
			ExpectContinueTimeout: 30 * time.Second,
			ResponseHeaderTimeout: 30 * time.Second,
		},
		Timeout: 30 * time.Second,
	}
}
