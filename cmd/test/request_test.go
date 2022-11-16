package test

import (
	"net/http"
	"testing"
)

func TestMakeRequestHeader(t *testing.T) {
	var (
		method        = http.MethodGet
		URL           = "http://localhost:9999"
		body   []byte = nil
	)

	testCases := []string{
		"1234",
		"4567",
	}

	for _, apiKey := range testCases {
		req, err := makeRequest(method, URL, body, apiKey)
		panicIf(err)

		if req.Header.Get("X-API-Key") != apiKey {
			t.Error("X-API-Key header is not specified properly")
		}
	}

}
