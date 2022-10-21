package main

import (
	"net/http"
	"testing"
)

func TestExtractEndpoint(t *testing.T) {
	testCases := []struct {
		name                 string
		method               string
		requestURL           string
		server               string
		expectedEndpointPath string
	}{
		{
			"relative url",
			http.MethodGet,
			"/driver_journeys",
			"",
			"/driver_journeys",
		},
		{
			"absolute url with trailing slash, api at root",
			http.MethodGet,
			"https://localhost:1323/driver_journeys",
			"https://localhost:1323/",
			"/driver_journeys",
		},
		{
			"absolute url without trailing slash, api at root",
			http.MethodGet,
			"https://localhost:1323/driver_journeys",
			"https://localhost:1323/",
			"/driver_journeys",
		},
		{
			"absolute url, api not at root",
			http.MethodGet,
			"https://localhost:1323/api/driver_journeys",
			"https://localhost:1323/api/",
			"/driver_journeys",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			request, err := http.NewRequest(tc.method, tc.requestURL, nil)
			panicIf(err)
			endpoint := ExtractEndpoint(request, tc.server)
			if endpoint.Method != tc.method || endpoint.Path != tc.expectedEndpointPath {
				t.Error("Failure to identify right endpoint from request")
			}
		})
	}
}
