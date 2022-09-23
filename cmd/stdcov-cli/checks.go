package main

import (
	"net/http"
	"net/url"
)

// CheckAPIStatus checks that `GET /status` endpoint returns a response with
// HTTP status "200 OK"
func CheckAPIStatus(client HTTPGetter, baseURL *url.URL) bool {
	statusURL := baseURL.JoinPath("/status")
	response, err := client.Get(statusURL.String())
	if err == nil && response.StatusCode == http.StatusOK {
		return true
	}
	return false
}
