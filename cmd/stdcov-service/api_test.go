package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	stdcovcli "gitlab.com/multi/stdcov-api-test/cmd/stdcov-cli"
	"gitlab.com/multi/stdcov-api-test/cmd/stdcov-service/server"
)

func TestCreateUser(t *testing.T) {
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
	a := stdcovcli.NewAssertionAccu()
	stdcovcli.TestGetDriverJourneys(request, response, a)
	for _, ar := range a.GetAssertionResults() {
		if err := ar.Unwrap(); err != nil {
			t.Log(err)
			t.Fail()
		}

	}
}
