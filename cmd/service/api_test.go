package service

import (
	"net/http"
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
		coordsIgnore = util.Coord{Lat: 0, Lon: 0}
		coordsRef    = util.Coord{Lat: 46.1604531, Lon: -1.2219607} // reference
		coords900m   = util.Coord{Lat: 46.1613442, Lon: -1.2103736} // at ~900m from reference
		coords1100m  = util.Coord{Lat: 46.1613679, Lon: -1.2086563} // at ~1100m from reference
		coords2100m  = util.Coord{Lat: 46.1649225, Lon: -1.1954497} // at ~2100m from reference
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
				makeDriverJourneyAtCoords(coords900m, coordsIgnore),
				makeDriverJourneyAtCoords(coords1100m, coordsIgnore),
			},
			false,
		},

		{
			"Departure radius 2",
			makeParamsWithDepartureRadius(coordsRef, 2),
			[]api.DriverJourney{
				makeDriverJourneyAtCoords(coords900m, coordsIgnore),
				makeDriverJourneyAtCoords(coords2100m, coordsIgnore),
			},
			false,
		},

		{
			"Departure radius 3",
			makeParamsWithDepartureRadius(coordsRef, 1),
			[]api.DriverJourney{
				makeDriverJourneyAtCoords(coords1100m, coordsIgnore),
			},
			true,
		},

		{
			"Departure radius 3",
			makeParamsWithDepartureRadius(coordsRef, 1),
			[]api.DriverJourney{
				makeDriverJourneyAtCoords(coords900m, coordsIgnore),
			},
			false,
		},

		{
			"Arrival radius 1",
			makeParamsWithArrivalRadius(coordsRef, 1),
			[]api.DriverJourney{
				makeDriverJourneyAtCoords(coordsIgnore, coords900m),
				makeDriverJourneyAtCoords(coordsIgnore, coords1100m),
			},
			false,
		},

		{
			"Arrival radius 2",
			makeParamsWithArrivalRadius(coordsRef, 2),
			[]api.DriverJourney{
				makeDriverJourneyAtCoords(coordsIgnore, coords2100m),
				makeDriverJourneyAtCoords(coordsIgnore, coords900m),
			},
			false,
		},

		{
			"Arrival radius 3",
			makeParamsWithArrivalRadius(coordsRef, 1),
			[]api.DriverJourney{
				makeDriverJourneyAtCoords(coordsIgnore, coords1100m),
			},
			true,
		},

		{
			"Arrival radius 4",
			makeParamsWithArrivalRadius(coordsRef, 1),
			[]api.DriverJourney{
				makeDriverJourneyAtCoords(coordsIgnore, coords900m),
			},
			false,
		},

		{
			"TimeDelta 1",
			makeParamsWithTimeDelta(10),
			[]api.DriverJourney{
				makeDriverJourneyAtDate(5),
			},
			false,
		},

		{
			"TimeDelta 2",
			makeParamsWithTimeDelta(10),
			[]api.DriverJourney{
				makeDriverJourneyAtDate(15),
			},
			true,
		},

		{
			"TimeDelta 3",
			makeParamsWithTimeDelta(20),
			[]api.DriverJourney{
				makeDriverJourneyAtDate(25),
				makeDriverJourneyAtDate(15),
			},
			false,
		},

		{
			"Count 1",
			makeParamsWithCount(1),
			makeNDriverJourneys(1),
			false,
		},

		{
			"Count 2",
			makeParamsWithCount(0),
			makeNDriverJourneys(1),
			true,
		},

		{
			"Count 3",
			makeParamsWithCount(2),
			makeNDriverJourneys(4),
			true,
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
	testRequest, err := params.MakeRequest(fakeServer)
	panicIf(err)

	mockDB := NewMockDB()
	mockDB.DriverJourneys = testData

	rec, ctx, handler := setupTest(testRequest, mockDB)

	// Make API Call
	err = handler.GetDriverJourneys(ctx, api.GetDriverJourneysParams(*params))
	panicIf(err)

	response := rec.Result()

	flags := test.Flags{DisallowEmpty: !expectEmpty}
	assertionResults := test.TestGetDriverJourneysResponse(testRequest, response, flags)
	checkAssertionResults(t, assertionResults)
}

func TestPassengerJourneys(t *testing.T) {

	var (
		coordsIgnore = util.Coord{Lat: 0, Lon: 0}
		coordsRef    = util.Coord{Lat: 46.1604531, Lon: -1.2219607} // reference
		coords900m   = util.Coord{Lat: 46.1613442, Lon: -1.2103736} // at ~900m from reference
		/* coords1100m  = util.Coord{Lat: 46.1613679, Lon: -1.2086563} // at ~1100m from reference */
	/* 	coords2100m  = util.Coord{Lat: 46.1649225, Lon: -1.1954497} // at ~2100m from reference */
	)

	testCases := []struct {
		name              string
		testParams        *api.GetPassengerJourneysParams
		testData          []api.PassengerJourney
		expectEmptyResult bool
	}{

		{"No data", &api.GetPassengerJourneysParams{}, []api.PassengerJourney{}, true},

		{
			"Departure radius 0",
			makeParamsWithDepartureRadius2(coordsRef, 1),
			[]api.PassengerJourney{
				makePassengerJourneyAtCoords(coords900m, coordsIgnore),
			},
			false,
		},

		/* { */
		/* 	"Departure radius 1", */
		/* 	makeParamsWithDepartureRadius2(coordsRef, 1), */
		/* 	[]api.PassengerJourney{ */
		/* 		makePassengerJourneyAtCoords(coords900m, coordsIgnore), */
		/* 		makePassengerJourneyAtCoords(coords1100m, coordsIgnore), */
		/* 	}, */
		/* 	false, */
		/* }, */

		/* { */
		/* 	"Departure radius 2", */
		/* 	makeParamsWithDepartureRadius(coordsRef, 2), */
		/* 	[]api.DriverJourney{ */
		/* 		makeDriverJourneyAtCoords(coords900m, coordsIgnore), */
		/* 		makeDriverJourneyAtCoords(coords2100m, coordsIgnore), */
		/* 	}, */
		/* 	false, */
		/* }, */

		/* { */
		/* 	"Departure radius 3", */
		/* 	makeParamsWithDepartureRadius(coordsRef, 1), */
		/* 	[]api.DriverJourney{ */
		/* 		makeDriverJourneyAtCoords(coords1100m, coordsIgnore), */
		/* 	}, */
		/* 	true, */
		/* }, */

		/* { */
		/* 	"Arrival radius 1", */
		/* 	makeParamsWithArrivalRadius(coordsRef, 1), */
		/* 	[]api.DriverJourney{ */
		/* 		makeDriverJourneyAtCoords(coordsIgnore, coords900m), */
		/* 		makeDriverJourneyAtCoords(coordsIgnore, coords1100m), */
		/* 	}, */
		/* 	false, */
		/* }, */

		/* { */
		/* 	"Arrival radius 2", */
		/* 	makeParamsWithArrivalRadius(coordsRef, 2), */
		/* 	[]api.DriverJourney{ */
		/* 		makeDriverJourneyAtCoords(coordsIgnore, coords2100m), */
		/* 		makeDriverJourneyAtCoords(coordsIgnore, coords900m), */
		/* 	}, */
		/* 	false, */
		/* }, */

		/* { */
		/* 	"Arrival radius 3", */
		/* 	makeParamsWithArrivalRadius(coordsRef, 1), */
		/* 	[]api.DriverJourney{ */
		/* 		makeDriverJourneyAtCoords(coordsIgnore, coords1100m), */
		/* 	}, */
		/* 	true, */
		/* }, */

		/* { */
		/* 	"Arrival radius 4", */
		/* 	makeParamsWithArrivalRadius(coordsRef, 1), */
		/* 	[]api.DriverJourney{ */
		/* 		makeDriverJourneyAtCoords(coordsIgnore, coords900m), */
		/* 	}, */
		/* 	false, */
		/* }, */

		/* { */
		/* 	"TimeDelta 1", */
		/* 	makeParamsWithTimeDelta(10), */
		/* 	[]api.DriverJourney{ */
		/* 		makeDriverJourneyAtDate(5), */
		/* 	}, */
		/* 	false, */
		/* }, */

		/* { */
		/* 	"TimeDelta 2", */
		/* 	makeParamsWithTimeDelta(10), */
		/* 	[]api.DriverJourney{ */
		/* 		makeDriverJourneyAtDate(15), */
		/* 	}, */
		/* 	true, */
		/* }, */

		/* { */
		/* 	"TimeDelta 3", */
		/* 	makeParamsWithTimeDelta(20), */
		/* 	[]api.DriverJourney{ */
		/* 		makeDriverJourneyAtDate(25), */
		/* 		makeDriverJourneyAtDate(15), */
		/* 	}, */
		/* 	false, */
		/* }, */

		/* { */
		/* 	"Count 1", */
		/* 	makeParamsWithCount(1), */
		/* 	makeNDriverJourneys(1), */
		/* 	false, */
		/* }, */

		/* { */
		/* 	"Count 2", */
		/* 	makeParamsWithCount(0), */
		/* 	makeNDriverJourneys(1), */
		/* 	true, */
		/* }, */

		/* { */
		/* 	"Count 3", */
		/* 	makeParamsWithCount(2), */
		/* 	makeNDriverJourneys(4), */
		/* 	true, */
		/* }, */
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			testGetPassengerJourneyRequestWithData(
				t,
				tc.testParams,
				tc.testData,
				tc.expectEmptyResult,
			)
		})
	}
}

func testGetPassengerJourneyRequestWithData(
	t *testing.T,
	params api.GetJourneysParams,
	testData []api.PassengerJourney,
	expectEmpty bool,
) {
	testRequest, err := params.MakeRequest(fakeServer)
	panicIf(err)

	mockDB := NewMockDB()
	mockDB.PassengerJourneys = testData

	rec, ctx, handler := setupTest(testRequest, mockDB)

	// Make API Call
	err = api.GetJourneys(handler, ctx, params)
	panicIf(err)

	response := rec.Result()
	flags := test.Flags{DisallowEmpty: !expectEmpty}
	assertionResults := test.TestGetPassengerJourneysResponse(testRequest, response, flags)
	checkAssertionResults(t, assertionResults)
}

func setupTest(request *http.Request, mockDB MockDB) (*httptest.ResponseRecorder, echo.Context, api.ServerInterface) {
	e := echo.New()
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(request, rec)
	handler := &StdCovServerImpl{mockDB}
	return rec, c, handler
}

func checkAssertionResults(t *testing.T, assertionResults []test.AssertionResult) {
	assert.Greater(t, len(assertionResults), 0)
	for _, ar := range assertionResults {
		if err := ar.Unwrap(); err != nil {
			t.Error(err)
		}
	}
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}
