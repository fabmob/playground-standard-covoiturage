package main

import (
	"errors"
	"net/http"
	"net/url"
	"testing"
)

var fakeBaseURL, _ = url.Parse("http://localhost:8000")

func TestApiStatus(t *testing.T) {
	t.Run("Status check with response", func(t *testing.T) {
		testCases := []struct {
			statusCode int
			isUp       bool
		}{
			{http.StatusOK, true},
			{http.StatusInternalServerError, false},
			{http.StatusTooManyRequests, false},
		}

		for _, tc := range testCases {
			m := returnStatusCodeClient(tc.statusCode)
			if CheckAPIStatus(m, fakeBaseURL) != tc.isUp {
				t.Errorf("Wrong behavior of status check with status %d", tc.statusCode)
			}
		}
	})

	t.Run("Status check with error", func(t *testing.T) {

		urlError := &url.Error{Op: "", URL: "", Err: errors.New("")}
		m := &mockClient{Error: urlError}
		if CheckAPIStatus(m, fakeBaseURL) {
			t.Error("If error returned, api is not up")
		}
	})

	t.Run("Status check hits the right endpoint", func(t *testing.T) {
		m := returnStatusCodeClient(http.StatusOK)
		_ = CheckAPIStatus(m, fakeBaseURL)
		if m.nCalls != 1 {
			t.Error("Exactly one request expected")
		}
		urlCalled, err := url.Parse(m.lastURLCalled)
		if err != nil {
			t.Error("CheckApiUp should request a valid url")
		}
		if urlCalled.Path != "/status" {
			t.Logf("Url: %s", urlCalled)
			t.Logf("Parsed path: %s", urlCalled.Path)
			t.Error("CheckApiUp should request the status endpoint")
		}
	})
}
