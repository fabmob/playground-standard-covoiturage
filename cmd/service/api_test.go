package service

import (
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gitlab.com/multi/stdcov-api-test/cmd/service/server"
	"gitlab.com/multi/stdcov-api-test/cmd/test"
	"gitlab.com/multi/stdcov-api-test/cmd/test/client"
)

var defaultTestFlags test.Flags = test.Flags{DisallowEmpty: false}

var fakeServer = "https:localhost:1323"

func TestDriverJourneys(t *testing.T) {

	var (
		coordsRef   = coords{46.1604531, -1.2219607} // reference
		coords900m  = coords{46.1613442, -1.2103736} // at ~900m from reference
		coords1100m = coords{46.1613679, -1.2086563} // at ~1100m from reference
	)

	testCases := []struct {
		name              string
		testParams        *client.GetDriverJourneysParams
		testData          []server.DriverJourney
		expectEmptyResult bool
	}{

		{"No data", &client.GetDriverJourneysParams{}, []server.DriverJourney{}, true},
		{
			"Departure radius",
			paramsWithDepartureRadius(coordsRef, 1),
			[]server.DriverJourney{
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
			},
			false,
		},
	}
	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			testGetDriverJourneyRequestWithData(t, tc.testParams, tc.testData)
		})
	}
}

func testGetDriverJourneyRequestWithData(
	t *testing.T,
	params *client.GetDriverJourneysParams,
	testData []server.DriverJourney,
) {

	testRequest, err := client.NewGetDriverJourneysRequest(fakeServer, params)
	panicIf(err)

	e := echo.New()
	testRequest.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(testRequest, rec)

	mockDB := NewMockDB()
	mockDB.driverJourneys = testData
	handler := &StdCovServerImpl{mockDB}

	// Assertions
	err = handler.GetDriverJourneys(c, server.GetDriverJourneysParams(*params))
	if err != nil {
		t.Fail()
	}
	response := rec.Result()
	a := test.NewAssertionAccu()
	test.TestGetDriverJourneysResponse(testRequest, response, a,
		defaultTestFlags)
	assert.Greater(t, len(a.GetAssertionResults()), 0)
	for _, ar := range a.GetAssertionResults() {
		if err := ar.Unwrap(); err != nil {
			t.Log(err)
			t.Fail()
		}
	}
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

type coords struct {
	lat float64
	lon float64
}

func paramsWithDepartureRadius(departureCoords coords, departureRadius float32) *client.GetDriverJourneysParams {
	params := client.NewGetDriverJourneysParams(
		float32(departureCoords.lat),
		float32(departureCoords.lon),
		0,
		0,
		0,
	)
	params.DepartureRadius = &departureRadius
	return params
}
