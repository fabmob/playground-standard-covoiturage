package test

import (
	"net/http"
	"testing"
)

func TestSplitServerEndpoint(t *testing.T) {
	testCases := []struct {
		name             string
		method           string
		requestURL       string
		expectedServer   string
		expectedEndpoint Endpoint
		expectError      bool
	}{
		{
			"simple case 1",
			http.MethodGet,
			"https://localhost:1323/passenger_journeys",
			"https://localhost:1323",
			GetPassengerJourneysEndpoint,
			false,
		},

		{
			"simple case 2",
			http.MethodGet,
			"https://localhost:1323/api/driver_journeys",
			"https://localhost:1323/api",
			GetDriverJourneysEndpoint,
			false,
		},

		{
			"wrong method",
			http.MethodPost,
			"https://localhost:1323/api/driver_journeys",
			"",
			Endpoint{},
			true,
		},

		{
			"more complex case 1: username & password",
			http.MethodGet,
			"http://username:password@example.com/a/b/c/driver_journeys",
			"http://username:password@example.com/a/b/c",
			GetDriverJourneysEndpoint,
			false,
		},

		{
			"more complex case 2: query",
			http.MethodGet,
			"http://example.com/a/b/c/driver_journeys?stuff=3",
			"http://example.com/a/b/c",
			GetDriverJourneysEndpoint,
			false,
		},

		{
			"more complex case 3: path parameter",
			http.MethodGet,
			"http://example.com/bookings/1234",
			"http://example.com",
			GetBookingsEndpoint,
			false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			guessedServer, guessedEndpoint, err := SplitServerEndpoint(tc.method, tc.requestURL)
			if tc.expectError != (err != nil) {
				t.Fail()
			}
			if guessedServer != tc.expectedServer {
				t.Logf("Expected server: %s", tc.expectedServer)
				t.Logf("Got: %s (error %s)", guessedServer, err)
				t.Fail()
			}
			if guessedEndpoint != tc.expectedEndpoint {
				t.Logf("Expected endpoint: %s", tc.expectedEndpoint)
				t.Logf("Got: %s (error %s)", guessedEndpoint, err)
				t.Fail()
			}
		})
	}
}
