package service

import (
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gitlab.com/multi/stdcov-api-test/cmd/api"
	"gitlab.com/multi/stdcov-api-test/cmd/test"
	"gitlab.com/multi/stdcov-api-test/cmd/util"
)

var fakeServer = "https:localhost:1323"

func TestDriverJourneys(t *testing.T) {

	var (
		coordsIgnore = util.Coord{0, 0}
		coordsRef    = util.Coord{46.1604531, -1.2219607} // reference
		coords900m   = util.Coord{46.1613442, -1.2103736} // at ~900m from reference
		coords1100m  = util.Coord{46.1613679, -1.2086563} // at ~1100m from reference
		coords2100m  = util.Coord{46.1649225, -1.1954497} // at ~2100m from reference
	)

	testCases := []struct {
		name              string
		testParams        *api.GetDriverJourneysParams
		testData          []api.DriverJourney
		expectEmptyResult bool
	}{

		{"No data", &api.GetDriverJourneysParams{}, []api.DriverJourney{}, true},

		{
			"Departure radius 1",
			makeParamsWithDepartureRadius(coordsRef, 1),
			[]api.DriverJourney{
				makeDriverJourney(coords900m, coordsIgnore),
				makeDriverJourney(coords1100m, coordsIgnore),
			},
			false,
		},

		{
			"Departure radius 2",
			makeParamsWithDepartureRadius(coordsRef, 2),
			[]api.DriverJourney{
				makeDriverJourney(coords900m, coordsIgnore),
				makeDriverJourney(coords2100m, coordsIgnore),
			},
			false,
		},

		{
			"Departure radius 4",
			makeParamsWithDepartureRadius(coordsRef, 1),
			[]api.DriverJourney{
				makeDriverJourney(coords1100m, coordsIgnore),
			},
			true,
		},

		{
			"Departure radius 3",
			makeParamsWithDepartureRadius(coordsRef, 1),
			[]api.DriverJourney{
				makeDriverJourney(coords900m, coordsIgnore),
			},
			false,
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			testGetDriverJourneyRequestWithData(
				t,
				tc.testParams,
				tc.testData,
				tc.expectEmptyResult,
			)
		})
	}
}

func testGetDriverJourneyRequestWithData(
	t *testing.T,
	params *api.GetDriverJourneysParams,
	testData []api.DriverJourney,
	expectEmpty bool,
) {

	testRequest, err := api.NewGetDriverJourneysRequest(fakeServer, params)
	panicIf(err)

	e := echo.New()
	testRequest.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(testRequest, rec)

	mockDB := NewMockDB()
	mockDB.driverJourneys = testData
	handler := &StdCovServerImpl{mockDB}

	// Assertions
	err = handler.GetDriverJourneys(c, api.GetDriverJourneysParams(*params))
	if err != nil {
		t.Fail()
	}
	response := rec.Result()
	a := test.NewAssertionAccu()
	flags := test.Flags{DisallowEmpty: !expectEmpty}
	test.TestGetDriverJourneysResponse(testRequest, response, a, flags)
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
