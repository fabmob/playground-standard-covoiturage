package test

import (
	"net/http"
	"testing"

	"github.com/fabmob/playground-standard-covoiturage/cmd/test/endpoint"
)

func TestSplitServerEndpoint(t *testing.T) {
	testCases := []struct {
		name             string
		method           string
		requestURL       string
		expectedServer   string
		expectedEndpoint endpoint.Info
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
			endpoint.Info{},
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

		{
			"more complex case 4: real example",
			http.MethodGet,
			"https://api-host.preprod-ab.some-domain.fr/api/path/1ab2c34-56d-343e21-f0g/other_stuff/driver_journeys?departureLat=48.8588548&departureLng=2.264463&arrivalLat=47.8733876&arrivalLng=1.8296428&departureDate=1668608335&timeDelta=100000&departureRadius=10&arrivalRadius=10",
			"https://api-host.preprod-ab.some-domain.fr/api/path/1ab2c34-56d-343e21-f0g/other_stuff",
			GetDriverJourneysEndpoint,
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
