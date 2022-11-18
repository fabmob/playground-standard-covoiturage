package service

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fabmob/playground-standard-covoiturage/cmd/api"
	"github.com/fabmob/playground-standard-covoiturage/cmd/test"
	"github.com/fabmob/playground-standard-covoiturage/cmd/util"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestDriverJourneys(t *testing.T) {
	var (
		coordsIgnore = util.Coord{Lat: 0, Lon: 0}
		coordsRef    = util.Coord{Lat: 46.1604531, Lon: -1.2219607} // reference
		coords900m   = util.Coord{Lat: 46.1613442, Lon: -1.2103736} // at ~900m from reference
		coords1100m  = util.Coord{Lat: 46.1613679, Lon: -1.2086563} // at ~1100m from reference
		coords2100m  = util.Coord{Lat: 46.1649225, Lon: -1.1954497} // at ~2100m from reference
	)

	testCases := []struct {
		name                 string
		testParams           api.GetJourneysParams
		testData             []api.DriverJourney
		expectNonEmptyResult bool
	}{

		{
			"No data",
			&api.GetDriverJourneysParams{},
			[]api.DriverJourney{},
			false,
		},

		{
			"Departure radius 1",
			makeParamsWithDepartureRadius(coordsRef, 1, "driver"),
			[]api.DriverJourney{
				makeDriverJourneyAtCoords(coords900m, coordsIgnore),
				makeDriverJourneyAtCoords(coords1100m, coordsIgnore),
			},
			true,
		},

		{
			"Departure radius 2",
			makeParamsWithDepartureRadius(coordsRef, 2, "driver"),
			[]api.DriverJourney{
				makeDriverJourneyAtCoords(coords900m, coordsIgnore),
				makeDriverJourneyAtCoords(coords2100m, coordsIgnore),
			},
			true,
		},

		{
			"Departure radius 3",
			makeParamsWithDepartureRadius(coordsRef, 1, "driver"),
			[]api.DriverJourney{
				makeDriverJourneyAtCoords(coords1100m, coordsIgnore),
			},
			false,
		},

		{
			"Departure radius 3",
			makeParamsWithDepartureRadius(coordsRef, 1, "driver"),
			[]api.DriverJourney{
				makeDriverJourneyAtCoords(coords900m, coordsIgnore),
			},
			true,
		},

		{
			"Arrival radius 1",
			makeParamsWithArrivalRadius(coordsRef, 1, "driver"),
			[]api.DriverJourney{
				makeDriverJourneyAtCoords(coordsIgnore, coords900m),
				makeDriverJourneyAtCoords(coordsIgnore, coords1100m),
			},
			true,
		},

		{
			"Arrival radius 2",
			makeParamsWithArrivalRadius(coordsRef, 2, "driver"),
			[]api.DriverJourney{
				makeDriverJourneyAtCoords(coordsIgnore, coords2100m),
				makeDriverJourneyAtCoords(coordsIgnore, coords900m),
			},
			true,
		},

		{
			"Arrival radius 3",
			makeParamsWithArrivalRadius(coordsRef, 1, "driver"),
			[]api.DriverJourney{
				makeDriverJourneyAtCoords(coordsIgnore, coords1100m),
			},
			false,
		},

		{
			"Arrival radius 4",
			makeParamsWithArrivalRadius(coordsRef, 1, "driver"),
			[]api.DriverJourney{
				makeDriverJourneyAtCoords(coordsIgnore, coords900m),
			},
			true,
		},

		{
			"TimeDelta 1",
			makeParamsWithTimeDelta(10, "driver"),
			[]api.DriverJourney{
				makeDriverJourneyAtDate(5),
			},
			true,
		},

		{
			"TimeDelta 2",
			makeParamsWithTimeDelta(10, "driver"),
			[]api.DriverJourney{
				makeDriverJourneyAtDate(15),
			},
			false,
		},

		{
			"TimeDelta 3",
			makeParamsWithTimeDelta(20, "driver"),
			[]api.DriverJourney{
				makeDriverJourneyAtDate(25),
				makeDriverJourneyAtDate(15),
			},
			true,
		},

		{
			"Count 1",
			makeParamsWithCount(1, "driver"),
			makeNDriverJourneys(1),
			true,
		},

		{
			"Count 2",
			makeParamsWithCount(0, "driver"),
			makeNDriverJourneys(1),
			false,
		},

		{
			"Count 3",
			makeParamsWithCount(2, "driver"),
			makeNDriverJourneys(4),
			true,
		},

		{
			"Count 4 - count > n driver journeys",
			makeParamsWithCount(1, "driver"),
			makeNDriverJourneys(0),
			false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testGetDriverJourneyRequestWithData(
				t,
				tc.testParams,
				tc.testData,
				tc.expectNonEmptyResult,
			)
		})
	}
}

func TestPassengerJourneys(t *testing.T) {
	var (
		coordsIgnore = util.Coord{Lat: 0, Lon: 0}
		coordsRef    = util.Coord{Lat: 46.1604531, Lon: -1.2219607} // reference
		coords900m   = util.Coord{Lat: 46.1613442, Lon: -1.2103736} // at ~900m from reference
		coords1100m  = util.Coord{Lat: 46.1613679, Lon: -1.2086563} // at ~1100m from reference
		coords2100m  = util.Coord{Lat: 46.1649225, Lon: -1.1954497} // at ~2100m from reference
	)

	testCases := []struct {
		name                 string
		testParams           api.GetJourneysParams
		testData             []api.PassengerJourney
		expectNonEmptyResult bool
	}{

		{
			"No data",
			&api.GetPassengerJourneysParams{},
			[]api.PassengerJourney{},
			false,
		},

		{
			"Departure radius 0",
			makeParamsWithDepartureRadius(coordsRef, 1, "passenger"),
			[]api.PassengerJourney{
				makePassengerJourneyAtCoords(coords900m, coordsIgnore),
			},
			true,
		},

		{
			"Departure radius 1",
			makeParamsWithDepartureRadius(coordsRef, 1, "passenger"),
			[]api.PassengerJourney{
				makePassengerJourneyAtCoords(coords900m, coordsIgnore),
				makePassengerJourneyAtCoords(coords1100m, coordsIgnore),
			},
			true,
		},

		{
			"Departure radius 2",
			makeParamsWithDepartureRadius(coordsRef, 2, "passenger"),
			[]api.PassengerJourney{
				makePassengerJourneyAtCoords(coords900m, coordsIgnore),
				makePassengerJourneyAtCoords(coords2100m, coordsIgnore),
			},
			true,
		},

		{
			"Departure radius 3",
			makeParamsWithDepartureRadius(coordsRef, 1, "passenger"),
			[]api.PassengerJourney{
				makePassengerJourneyAtCoords(coords1100m, coordsIgnore),
			},
			false,
		},

		{
			"Arrival radius 1",
			makeParamsWithArrivalRadius(coordsRef, 1, "passenger"),
			[]api.PassengerJourney{
				makePassengerJourneyAtCoords(coordsIgnore, coords900m),
				makePassengerJourneyAtCoords(coordsIgnore, coords1100m),
			},
			true,
		},

		{
			"Arrival radius 2",
			makeParamsWithArrivalRadius(coordsRef, 2, "passenger"),
			[]api.PassengerJourney{
				makePassengerJourneyAtCoords(coordsIgnore, coords2100m),
				makePassengerJourneyAtCoords(coordsIgnore, coords900m),
			},
			true,
		},

		{
			"Arrival radius 3",
			makeParamsWithArrivalRadius(coordsRef, 1, "passenger"),
			[]api.PassengerJourney{
				makePassengerJourneyAtCoords(coordsIgnore, coords1100m),
			},
			false,
		},

		{
			"Arrival radius 4",
			makeParamsWithArrivalRadius(coordsRef, 1, "passenger"),
			[]api.PassengerJourney{
				makePassengerJourneyAtCoords(coordsIgnore, coords900m),
			},
			true,
		},

		{
			"TimeDelta 1",
			makeParamsWithTimeDelta(10, "passenger"),
			[]api.PassengerJourney{
				makePassengerJourneyAtDate(5),
			},
			true,
		},

		{
			"TimeDelta 2",
			makeParamsWithTimeDelta(10, "passenger"),
			[]api.PassengerJourney{
				makePassengerJourneyAtDate(15),
			},
			false,
		},

		{
			"TimeDelta 3",
			makeParamsWithTimeDelta(20, "passenger"),
			[]api.PassengerJourney{
				makePassengerJourneyAtDate(25),
				makePassengerJourneyAtDate(15),
			},
			true,
		},

		{
			"Count 1",
			makeParamsWithCount(1, "passenger"),
			makeNPassengerJourneys(1),
			true,
		},

		{
			"Count 2",
			makeParamsWithCount(0, "passenger"),
			makeNPassengerJourneys(1),
			false,
		},

		{
			"Count 3",
			makeParamsWithCount(2, "passenger"),
			makeNPassengerJourneys(4),
			true,
		},

		{
			"Count 4 - count > n passenger journeys",
			makeParamsWithCount(1, "passenger"),
			makeNPassengerJourneys(0),
			false,
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			testGetPassengerJourneyRequestWithData(
				t,
				tc.testParams,
				tc.testData,
				tc.expectNonEmptyResult,
			)
		})
	}
}

func TestGetBookings(t *testing.T) {

	testCases := []struct {
		bookings           []api.Booking
		queryBookingID     uuid.UUID
		disallowEmpty      bool
		expectedStatusCode int
	}{
		{
			[]api.Booking{},
			repUUID(1),
			false,
			http.StatusNotFound,
		},
		{
			[]api.Booking{makeBooking(repUUID(2))},
			repUUID(2),
			true,
			http.StatusOK,
		},
	}

	for _, tc := range testCases {

		t.Run("test case", func(t *testing.T) {
			mockDB := NewMockDB()
			mockDB.Bookings = tc.bookings

			request, err := api.NewGetBookingsRequest(fakeServer, tc.queryBookingID)
			panicIf(err)

			// Setup testing server with response recorder
			e := echo.New()
			rec := httptest.NewRecorder()
			ctx := e.NewContext(request, rec)
			handler := NewServerWithDB(mockDB)

			// Make API Call
			err = handler.GetBookings(ctx, tc.queryBookingID)
			panicIf(err)

			response := rec.Result()
			flags := test.NewFlags()
			flags.DisallowEmpty = tc.disallowEmpty
			flags.ExpectedStatusCode = tc.expectedStatusCode

			assertionResults := test.TestGetBookingsResponse(request, response, flags)

			checkAssertionResults(t, assertionResults)
		})
	}
}

func TestPostBookings(t *testing.T) {

	request, err := api.NewPostBookingsRequest(fakeServer, api.Booking{})
	panicIf(err)

	// Setup testing server with response recorder
	e := echo.New()
	rec := httptest.NewRecorder()
	ctx := e.NewContext(request, rec)
	handler := NewServer()

	// Make API Call
	err = handler.PostBookings(ctx)
	panicIf(err)

	response := rec.Result()
	flags := test.NewFlags()
	flags.ExpectedStatusCode = http.StatusCreated

	assertionResults := test.TestGetBookingsResponse(request, response, flags)

	checkAssertionResults(t, assertionResults)
}

func testGetDriverJourneyRequestWithData(
	t *testing.T,
	params api.GetJourneysParams,
	testData []api.DriverJourney,
	expectNonEmpty bool,
) {

	mockDB := NewMockDB()
	mockDB.DriverJourneys = testData
	testFunction := test.TestGetDriverJourneysResponse

	testGetJourneys(t, params, mockDB, testFunction, expectNonEmpty)
}

func testGetPassengerJourneyRequestWithData(
	t *testing.T,
	params api.GetJourneysParams,
	testData []api.PassengerJourney,
	expectNonEmpty bool,
) {

	mockDB := NewMockDB()
	mockDB.PassengerJourneys = testData
	testFunction := test.TestGetPassengerJourneysResponse

	testGetJourneys(t, params, mockDB, testFunction, expectNonEmpty)
}

func testGetJourneys(t *testing.T, params api.GetJourneysParams, mockDB *MockDB, f test.ResponseTestFun, expectNonEmpty bool) {
	t.Helper()

	// Build request
	request, err := params.MakeRequest(fakeServer)
	panicIf(err)
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Setup testing server with response recorder
	e := echo.New()
	rec := httptest.NewRecorder()
	ctx := e.NewContext(request, rec)
	handler := NewServerWithDB(mockDB)

	// Make API Call
	err = api.GetJourneys(handler, ctx, params)
	panicIf(err)

	// Check response
	response := rec.Result()
	flags := test.NewFlags()
	flags.DisallowEmpty = expectNonEmpty
	assertionResults := f(request, response, flags)

	checkAssertionResults(t, assertionResults)
}

func checkAssertionResults(t *testing.T, assertionResults []test.AssertionResult) {
	t.Helper()

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
