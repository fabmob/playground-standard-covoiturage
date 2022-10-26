package service

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gitlab.com/multi/stdcov-api-test/cmd/service/server"
	"gitlab.com/multi/stdcov-api-test/cmd/test"
)

var defaultTestFlags test.Flags = test.Flags{DisallowEmpty: false}

func TestCreateUser(t *testing.T) {
	t.Run("", func(t *testing.T) {
		// Setup
		e := echo.New()
		request := httptest.NewRequest(
			http.MethodGet,
			"https://fabmob.github.io/driver_journeys?departureLat=0&departureLng=0&arrivalLat=0&arrivalLng=0&departureDate=0&timeDelta=900&departureRadius=1&arrivalRadius=1",
			nil,
		)
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(request, rec)

		mockDB := NewMockDB()
		mockDB.driverJourneys = []server.DriverJourney{}
		handler := &StdCovServerImpl{mockDB}

		params := server.GetDriverJourneysParams{}

		// Assertions
		err := handler.GetDriverJourneys(c, params)
		if err != nil {
			t.Fail()
		}
		response := rec.Result()
		a := test.NewAssertionAccu()
		a.Run(test.TestGetDriverJourneys(request, response, a,
			defaultTestFlags)...)
		assert.Greater(t, len(a.GetAssertionResults()), 0)
		for _, ar := range a.GetAssertionResults() {
			if err := ar.Unwrap(); err != nil {
				t.Log(err)
				t.Fail()
			}

		}
	})

	t.Run("DepartureRadius", func(t *testing.T) {
		// Setup
		e := echo.New()
		type coords struct {
			lat float64
			lon float64
		}
		var (
			/* coordsRef   = coords{46.1604531, -1.2219607} // reference */
			coords900m  = coords{46.1613442, -1.2103736} // at ~900m from reference
			coords1100m = coords{46.1613679, -1.2086563} // at ~1100m from reference
		)
		request := httptest.NewRequest(
			http.MethodGet,
			"https://fabmob.github.io/driver_journeys?departureLat=46.1604531&departureLng=-1.2219607&arrivalLat=0&arrivalLng=0&departureDate=0&timeDelta=900&departureRadius=1&arrivalRadius=1",
			nil,
		)
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(request, rec)

		mockDB := NewMockDB()
		mockDB.driverJourneys = []server.DriverJourney{
			{
				PassengerPickupLat: coords900m.lat,
				PassengerPickupLng: coords900m.lon,
				Type:               "DYNAMIC",
			},
			{
				PassengerPickupLat: coords1100m.lat,
				PassengerPickupLng: coords1100m.lon,
				Type:               "DYNAMIC",
			},
		}
		handler := &StdCovServerImpl{mockDB}

		params := server.GetDriverJourneysParams{}

		// Assertions
		err := handler.GetDriverJourneys(c, params)
		if err != nil {
			t.Fail()
		}
		response := rec.Result()
		a := test.NewAssertionAccu()
		a.Run(test.TestGetDriverJourneys(request, response, a,
			defaultTestFlags)...)
		assert.Greater(t, len(a.GetAssertionResults()), 0)
		for _, ar := range a.GetAssertionResults() {
			if err := ar.Unwrap(); err != nil {
				t.Log(err)
				t.Fail()
			}
		}
	})
}
