package service

import (
	"fmt"
	"net/http"
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

			mockDB := NewMockDB()
			mockDB.DriverJourneys = tc.testData

			flags := test.NewFlags()
			flags.DisallowEmpty = tc.expectNonEmptyResult

			testGetDriverJourneyHelper(
				t,
				tc.testParams,
				mockDB,
				flags,
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
			mockDB := NewMockDB()
			mockDB.PassengerJourneys = tc.testData

			flags := test.NewFlags()
			flags.DisallowEmpty = tc.expectNonEmptyResult

			testGetPassengerJourneyHelper(
				t,
				tc.testParams,
				mockDB,
				flags,
			)
		})
	}
}

func TestGetBookings(t *testing.T) {

	testCases := []struct {
		bookings           BookingsByID
		queryBookingID     uuid.UUID
		disallowEmpty      bool
		expectedStatusCode int
	}{
		{
			NewBookingsByID(),
			repUUID(1),
			false,
			http.StatusNotFound,
		},
		{
			NewBookingsByID(
				makeBooking(repUUID(2)),
			),
			repUUID(2),
			true,
			http.StatusOK,
		},
		{
			NewBookingsByID(
				makeBooking(repUUID(3)),
				makeBooking(repUUID(4)),
			),
			repUUID(4),
			true,
			http.StatusOK,
		},
	}

	for _, tc := range testCases {

		t.Run("test case", func(t *testing.T) {
			mockDB := NewMockDB()
			mockDB.Bookings = tc.bookings

			flags := test.NewFlags()
			flags.DisallowEmpty = tc.disallowEmpty
			flags.ExpectedStatusCode = tc.expectedStatusCode

			testGetBookingsHelper(t, mockDB, tc.queryBookingID, flags)
		})
	}
}

func TestPostBookings(t *testing.T) {

	testCases := []struct {
		booking              *api.Booking
		existingBookings     BookingsByID
		expectPostStatusCode int
		expectGetNonEmpty    bool
	}{
		{
			makeBooking(repUUID(10)),
			NewBookingsByID(),
			http.StatusCreated,
			true,
		},

		{
			makeBooking(repUUID(11)),
			NewBookingsByID(
				makeBooking(repUUID(11)),
			),
			http.StatusBadRequest,
			true,
		},
	}

	for _, tc := range testCases {
		bookingID := tc.booking.Id

		mockDB := NewMockDB()
		mockDB.Bookings = tc.existingBookings

		flagsPost := test.NewFlags()
		flagsPost.ExpectedStatusCode = tc.expectPostStatusCode

		flagsGet := test.NewFlags()
		flagsGet.DisallowEmpty = tc.expectGetNonEmpty

		t.Log(*tc.booking)
		testPostBookingsHelper(t, mockDB, *tc.booking, flagsPost)

		testGetBookingsHelper(t, mockDB, bookingID, flagsGet)
	}
}

func TestPatchBookings(t *testing.T) {

	testCases := []struct {
		bookingID               uuid.UUID
		newStatus               api.BookingStatus
		existingBookings        BookingsByID
		expectedPatchStatusCode int
		expectedGetStatusCode   int
		expectedStatus          api.BookingStatus
	}{
		{
			repUUID(20),
			api.BookingStatusVALIDATED,
			NewBookingsByID(
				makeBooking(repUUID(20)),
			),
			200,
			200,
			api.BookingStatusVALIDATED,
		},

		{
			repUUID(21),
			api.BookingStatusCOMPLETEDPENDINGVALIDATION,
			NewBookingsByID(
				makeBooking(repUUID(21)),
			),
			200,
			200,
			api.BookingStatusCOMPLETEDPENDINGVALIDATION,
		},

		{
			repUUID(22),
			api.BookingStatusCANCELLED,
			NewBookingsByID(),
			404,
			404,
			"",
		},

		{
			repUUID(23),
			api.BookingStatusVALIDATED,
			NewBookingsByID(
				makeBookingWithStatus(repUUID(23), api.BookingStatusVALIDATED),
			),
			409,
			200,
			api.BookingStatusVALIDATED,
		},

		{
			repUUID(24),
			api.BookingStatusVALIDATED,
			NewBookingsByID(
				makeBookingWithStatus(repUUID(24),
					api.BookingStatusCANCELLED),
			),
			409,
			200,
			api.BookingStatusCANCELLED,
		},

		{
			repUUID(25),
			"INVALID_STATUS",
			NewBookingsByID(
				makeBooking(repUUID(25)),
			),
			400,
			200,
			api.BookingStatusWAITINGCONFIRMATION,
		},
	}

	for _, tc := range testCases {

		mockDB := NewMockDB()
		mockDB.Bookings = tc.existingBookings

		params := api.PatchBookingsParams{Status: tc.newStatus}

		flagsPatch := test.NewFlags()
		flagsPatch.ExpectedStatusCode = tc.expectedPatchStatusCode

		flagsGet := test.NewFlags()
		flagsGet.ExpectedStatusCode = tc.expectedGetStatusCode
		flagsGet.ExpectedBookingStatus = tc.expectedStatus

		testPatchBookingsHelper(t, mockDB, tc.bookingID, params, flagsPatch)

		testGetBookingsHelper(t, mockDB, tc.bookingID, flagsGet)
	}
}

func TestPostBookingEvents(t *testing.T) {

	testCases := []struct {
		bookingID              uuid.UUID
		carpoolBookingEvent    *api.CarpoolBookingEvent
		existingBookings       BookingsByID
		expectedPostStatusCode int
		expectedGetStatusCode  int
		expectedStatus         api.BookingStatus
	}{
		{
			repUUID(31),
			makeCarpoolBookingEvent(repUUID(30), repUUID(31)),
			NewBookingsByID(),
			http.StatusOK,
			http.StatusOK,
			api.BookingStatusWAITINGCONFIRMATION,
		},

		{
			repUUID(32),
			makeCarpoolBookingEventWithStatus(repUUID(33), repUUID(32), api.BookingStatusCONFIRMED),
			NewBookingsByID(makeBooking(repUUID(32))),
			http.StatusOK,
			http.StatusOK,
			api.BookingStatusWAITINGCONFIRMATION,
		},

		{
			repUUID(34),
			makeCarpoolBookingEventWithStatus(repUUID(35), repUUID(34), api.BookingStatusCONFIRMED),
			NewBookingsByID(makeBookingWithStatus(repUUID(34), api.BookingStatusCONFIRMED)),
			http.StatusBadRequest,
			http.StatusOK,
			api.BookingStatusWAITINGCONFIRMATION,
		},

		{
			repUUID(36),
			makeCarpoolBookingEventWithStatus(repUUID(37), repUUID(36), api.BookingStatusCONFIRMED),
			NewBookingsByID(makeBookingWithStatus(repUUID(36), api.BookingStatusCANCELLED)),
			http.StatusBadRequest,
			http.StatusOK,
			api.BookingStatusWAITINGCONFIRMATION,
		},
	}

	for _, tc := range testCases {

		mockDB := NewMockDB()
		mockDB.Bookings = tc.existingBookings

		flagsPost := test.NewFlags()
		flagsPost.ExpectedStatusCode = tc.expectedPostStatusCode

		flagsGet := test.NewFlags()
		flagsGet.ExpectedStatusCode = tc.expectedGetStatusCode

		testPostBookingEventsHelper(t, mockDB, tc.carpoolBookingEvent, flagsPost)

		testGetBookingsHelper(t, mockDB, tc.bookingID, flagsGet)
	}
}

func TestPostMessage(t *testing.T) {
	testCases := []struct {
		message            api.PostMessagesJSONBody
		expectedStatusCode int
	}{
		{
			makeMessage(api.User{Id: "1", Alias: "quidam1"}, api.User{Id: "2",
				Alias: "quidam2"}),
			http.StatusCreated,
		},
	}
	for _, tc := range testCases {
		mockDB := NewMockDB()

		flags := test.NewFlags()
		flags.ExpectedStatusCode = tc.expectedStatusCode

		testPostMessageHelper(t, mockDB, tc.message, flags)

	}
}

func testPostMessageHelper(t *testing.T, mockDB *MockDB, message api.PostMessagesJSONBody, flags test.Flags) {
	request, err := api.NewPostMessagesRequest(fakeServer,
		api.PostMessagesJSONRequestBody(message))
	panicIf(err)

	// Setup testing server with response recorder
	handler, ctx, rec := setupTestServer(mockDB, request)

	// Make API Call
	err = handler.PostMessages(ctx)
	panicIf(err)

	// Test response
	response := rec.Result()
	fmt.Println(response)

	assertionResults := test.TestPostMessagesResponse(request, response, flags)
	checkAssertionResults(t, assertionResults)
}

func testPostBookingEventsHelper(t *testing.T, mockDB *MockDB,
	carpoolBookingEvent *api.CarpoolBookingEvent, flags test.Flags) {

	request, err := api.NewPostBookingEventsRequest(fakeServer, *carpoolBookingEvent)
	panicIf(err)

	// Setup testing server with response recorder
	handler, ctx, rec := setupTestServer(mockDB, request)

	// Make API Call
	err = handler.PostBookingEvents(ctx)
	panicIf(err)

	response := rec.Result()

	assertionResults := test.TestPostBookingEventsResponse(request, response, flags)
	checkAssertionResults(t, assertionResults)
}

func testPostBookingsHelper(
	t *testing.T,
	mockDB *MockDB,
	booking api.Booking,
	flags test.Flags,
) {
	t.Helper()

	request, err := api.NewPostBookingsRequest(fakeServer, booking)
	panicIf(err)

	// Setup testing server with response recorder
	handler, ctx, rec := setupTestServer(mockDB, request)

	// Make API Call
	err = handler.PostBookings(ctx)
	panicIf(err)

	response := rec.Result()

	assertionResults := test.TestPostBookingsResponse(request, response, flags)
	checkAssertionResults(t, assertionResults)
}

func testGetBookingsHelper(
	t *testing.T,
	mockDB *MockDB,
	bookingID api.BookingId,
	flags test.Flags,
) {
	t.Helper()

	// Make Request
	request, err := api.NewGetBookingsRequest(fakeServer, bookingID)
	panicIf(err)

	// Setup testing server with response recorder
	handler, ctx, rec := setupTestServer(mockDB, request)

	// Make API call
	err = handler.GetBookings(ctx, bookingID)
	panicIf(err)

	// Test results
	response := rec.Result()

	assertionResults := test.TestGetBookingsResponse(request, response, flags)

	checkAssertionResults(t, assertionResults)
}

func testPatchBookingsHelper(
	t *testing.T,
	mockDB *MockDB,
	bookingID api.BookingId,
	params api.PatchBookingsParams,
	flags test.Flags,
) {
	t.Helper()

	// Make Request
	request, err := api.NewPatchBookingsRequest(fakeServer, bookingID, &params)
	panicIf(err)

	// Setup testing server with response recorder
	handler, ctx, rec := setupTestServer(mockDB, request)

	// Make API call
	err = handler.PatchBookings(ctx, bookingID, params)
	panicIf(err)

	// Test results
	response := rec.Result()

	assertionResults := test.TestPatchBookingsResponse(request, response, flags)

	checkAssertionResults(t, assertionResults)
}

func testGetDriverJourneyHelper(
	t *testing.T,
	params api.GetJourneysParams,
	mockDB *MockDB,
	flags test.Flags,
) {

	testFunction := test.TestGetDriverJourneysResponse

	testGetJourneysHelper(t, params, mockDB, testFunction, flags)
}

func testGetPassengerJourneyHelper(
	t *testing.T,
	params api.GetJourneysParams,
	mockDB *MockDB,
	flags test.Flags,
) {

	testFunction := test.TestGetPassengerJourneysResponse

	testGetJourneysHelper(t, params, mockDB, testFunction, flags)
}

func testGetJourneysHelper(t *testing.T, params api.GetJourneysParams, mockDB *MockDB, f test.ResponseTestFun, flags test.Flags) {
	t.Helper()

	// Build request
	request, err := params.MakeRequest(fakeServer)
	panicIf(err)
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Setup testing server with response recorder
	handler, ctx, rec := setupTestServer(mockDB, request)

	// Make API Call
	err = api.GetJourneys(handler, ctx, params)
	panicIf(err)

	// Check response
	response := rec.Result()
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
