package test

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
		{
			"/get bookings endpoint",
			http.MethodGet,
			"https://localhost:1323/bookings/1234",
			"https://localhost:1323/api/",
			"/bookings",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			request, err := http.NewRequest(tc.method, tc.requestURL, nil)
			panicIf(err)
			endpoint, _ := ExtractEndpoint(request, tc.server)
			if endpoint.Method != tc.method || endpoint.Path != tc.expectedEndpointPath {
				t.Logf("Method : exected %s, got %s", tc.method, endpoint.Method)
				t.Logf("Path : exected %s, got %s", tc.expectedEndpointPath, endpoint.Path)
				t.Error("Failure to identify right endpoint from request")
			}
		})
	}
}

func TestGuessServer(t *testing.T) {
	testCases := []struct {
		name           string
		method         string
		requestURL     string
		expectedServer string
		expectError    bool
	}{
		{
			"simple case 1",
			http.MethodGet,
			"https://localhost:1323/passenger_journeys",
			"https://localhost:1323",
			false,
		},

		{
			"simple case 2",
			http.MethodGet,
			"https://localhost:1323/api/driver_journeys",
			"https://localhost:1323/api",
			false,
		},

		{
			"wrong method",
			http.MethodPost,
			"https://localhost:1323/api/driver_journeys",
			"",
			true,
		},

		{
			"more complex case 1: username & password",
			http.MethodGet,
			"http://username:password@example.com/a/b/c/driver_journeys",
			"http://username:password@example.com/a/b/c",
			false,
		},

		{
			"more complex case 2: query",
			http.MethodGet,
			"http://example.com/a/b/c/driver_journeys?stuff=3",
			"http://example.com/a/b/c",
			false,
		},

		{
			"more complex case 3: path parameter",
			http.MethodGet,
			"http://example.com/bookings/1234",
			"http://example.com",
			false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			guessedServer, err := GuessServer(tc.method, tc.requestURL)
			if tc.expectError != (err != nil) {
				t.Fail()
			}
			if guessedServer != tc.expectedServer {
				t.Logf("Expected server: %s", tc.expectedServer)
				t.Logf("Got: %s (error %s)", guessedServer, err)
				t.Fail()
			}
		})
	}
}
