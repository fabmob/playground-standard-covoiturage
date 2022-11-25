package service

import (
	"bytes"
	"flag"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/fabmob/playground-standard-covoiturage/cmd/api"
	"github.com/fabmob/playground-standard-covoiturage/cmd/endpoint"
	"github.com/fabmob/playground-standard-covoiturage/cmd/service/db"
	"github.com/fabmob/playground-standard-covoiturage/cmd/test"
	"github.com/fabmob/playground-standard-covoiturage/cmd/util"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func init() {
	// test flags do not need to be parsed explicitely, as it is already done in
	// normal `go test` operation
	flag.BoolVar(&generateTestData, "generate", false, "Should test data be regenerated")
}

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

			mockDB := db.NewMockDB()
			mockDB.DriverJourneys = tc.testData

			flags := test.NewFlags()
			flags.ExpectNonEmpty = tc.expectNonEmptyResult

			testGetDriverJourneyHelper(
				t,
				tc.testParams,
				mockDB,
				flags,
			)
		})
	}
}

func testGetDriverJourneyHelper(
	t *testing.T,
	params api.GetJourneysParams,
	mockDB *db.Mock,
	flags test.Flags,
) {
	testFunction := test.TestGetDriverJourneysResponse

	testGetJourneysHelper(t, params, mockDB, testFunction, flags)
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
			mockDB := db.NewMockDB()
			mockDB.PassengerJourneys = tc.testData

			flags := test.NewFlags()
			flags.ExpectNonEmpty = tc.expectNonEmptyResult

			testGetPassengerJourneyHelper(
				t,
				tc.testParams,
				mockDB,
				flags,
			)
		})
	}
}

func testGetPassengerJourneyHelper(
	t *testing.T,
	params api.GetJourneysParams,
	mockDB *db.Mock,
	flags test.Flags,
) {
	testFunction := test.TestGetPassengerJourneysResponse

	testGetJourneysHelper(t, params, mockDB, testFunction, flags)
}

func testGetJourneysHelper(t *testing.T, params api.GetJourneysParams, mockDB *db.Mock, f test.ResponseTestFun, flags test.Flags) {
	t.Helper()

	appendDataIfGenerated(t, mockDB)

	// Build request
	request, err := params.MakeRequest(localServer)
	util.PanicIf(err)
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request, err = endpoint.AddEndpointContext(request)
	util.PanicIf(err)

	// Setup testing server with response recorder
	handler, ctx, rec := setupTestServer(mockDB, request)

	// Make API Call
	err = api.GetJourneys(handler, ctx, params)
	util.PanicIf(err)

	// Check response
	response := rec.Result()
	assertionResults := f(request, response, flags)

	checkAssertionResults(t, assertionResults)
}

func TestGetBookings(t *testing.T) {
	testCases := []struct {
		name               string
		bookings           db.BookingsByID
		queryBookingID     uuid.UUID
		disallowEmpty      bool
		expectedStatusCode int
	}{
		{
			"getting a non-existing booking returns code 404",
			NewBookingsByID(),
			repUUID(1),
			false,
			http.StatusNotFound,
		},
		{
			"getting an existing booking returns it with code 200 #1",
			NewBookingsByID(
				makeBooking(repUUID(2)),
			),
			repUUID(2),
			true,
			http.StatusOK,
		},
		{
			"getting an existing booking returns it with code 200 #2",
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

		t.Run(tc.name, func(t *testing.T) {
			mockDB := db.NewMockDB()
			mockDB.Bookings = tc.bookings

			flags := test.NewFlags()
			flags.ExpectNonEmpty = tc.disallowEmpty
			flags.ExpectedResponseCode = tc.expectedStatusCode

			TestGetBookingsHelper(t, mockDB, tc.queryBookingID, flags)
		})
	}
}

func TestPostBookings(t *testing.T) {

	testCases := []struct {
		name                 string
		booking              *api.Booking
		existingBookings     db.BookingsByID
		expectPostStatusCode int
		expectGetNonEmpty    bool
	}{
		{
			"Posting a new booking succeeds with code 201",
			makeBooking(repUUID(10)),
			NewBookingsByID(),
			http.StatusCreated,
			true,
		},

		{
			"Posting a booking with colliding ID fails with code 400",
			makeBooking(repUUID(11)),
			NewBookingsByID(
				makeBooking(repUUID(11)),
			),
			http.StatusBadRequest,
			true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bookingID := tc.booking.Id

			mockDB := db.NewMockDB()
			mockDB.Bookings = tc.existingBookings

			flagsPost := test.NewFlags()
			flagsPost.ExpectedResponseCode = tc.expectPostStatusCode

			flagsGet := test.NewFlags()
			flagsGet.ExpectNonEmpty = tc.expectGetNonEmpty

			TestPostBookingsHelper(t, mockDB, *tc.booking, flagsPost)

			TestGetBookingsHelper(t, mockDB, bookingID, flagsGet)
		})
	}
}

func TestPatchBookings(t *testing.T) {

	testCases := []struct {
		name                    string
		bookingID               uuid.UUID
		newStatus               api.BookingStatus
		existingBookings        db.BookingsByID
		expectedPatchStatusCode int
		expectedGetStatusCode   int
		expectedStatus          api.BookingStatus
	}{
		{
			"patching VALIDATED over WAITING_CONFIRMATION succeeds",
			repUUID(20),
			api.BookingStatusVALIDATED,
			NewBookingsByID(
				makeBooking(repUUID(20)),
			),
			http.StatusOK,
			http.StatusOK,
			api.BookingStatusVALIDATED,
		},

		{
			"patching COMPLETED_PENDING_VALIDATION over WAITING_CONFIRMATION succeeds",
			repUUID(21),
			api.BookingStatusCOMPLETEDPENDINGVALIDATION,
			NewBookingsByID(
				makeBooking(repUUID(21)),
			),
			http.StatusOK,
			http.StatusOK,
			api.BookingStatusCOMPLETEDPENDINGVALIDATION,
		},

		{
			"patching a non-existing booking returns code 404",
			repUUID(22),
			api.BookingStatusCANCELLED,
			NewBookingsByID(),
			http.StatusNotFound,
			http.StatusNotFound,
			"",
		},

		{
			"patching VALIDATED other VALIDATED fails with code 409",
			repUUID(23),
			api.BookingStatusVALIDATED,
			NewBookingsByID(
				makeBookingWithStatus(repUUID(23), api.BookingStatusVALIDATED),
			),
			http.StatusConflict,
			http.StatusOK,
			api.BookingStatusVALIDATED,
		},

		{
			"patching VALIDATED other CANCELLED fails with code 409",
			repUUID(24),
			api.BookingStatusVALIDATED,
			NewBookingsByID(
				makeBookingWithStatus(repUUID(24),
					api.BookingStatusCANCELLED),
			),
			http.StatusConflict,
			http.StatusOK,
			api.BookingStatusCANCELLED,
		},

		{
			"patching INVALID_STATUS fails with code 400",
			repUUID(25),
			"INVALID_STATUS",
			NewBookingsByID(
				makeBooking(repUUID(25)),
			),
			http.StatusBadRequest,
			http.StatusOK,
			api.BookingStatusWAITINGCONFIRMATION,
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			mockDB := db.NewMockDB()
			mockDB.Bookings = tc.existingBookings

			flagsPatch := test.NewFlags()
			flagsPatch.ExpectedResponseCode = tc.expectedPatchStatusCode

			flagsGet := test.NewFlags()
			flagsGet.ExpectedResponseCode = tc.expectedGetStatusCode
			flagsGet.ExpectedBookingStatus = tc.expectedStatus

			TestPatchBookingsHelper(t, mockDB, tc.bookingID, tc.newStatus, flagsPatch)

			TestGetBookingsHelper(t, mockDB, tc.bookingID, flagsGet)
		})
	}
}

func TestPostBookingEvents(t *testing.T) {

	testCases := []struct {
		name                   string
		bookingID              uuid.UUID
		carpoolBookingEvent    *api.CarpoolBookingEvent
		existingBookings       db.BookingsByID
		expectedPostStatusCode int
		expectedGetStatusCode  int
		expectedBookingStatus  api.BookingStatus
	}{
		{
			"posting a new bookingEvent with status WAITING_CONFIRMATION succeeds",
			repUUID(31),
			makeCarpoolBookingEvent(repUUID(30), repUUID(31)),
			NewBookingsByID(),
			http.StatusOK,
			http.StatusOK,
			api.BookingStatusWAITINGCONFIRMATION,
		},

		{
			"posting a bookingEvent on existing booking (status CONFIRMED over WAITING_CONFIRMATION) changes its status",
			repUUID(32),
			makeCarpoolBookingEventWithStatus(repUUID(33), repUUID(32), api.BookingStatusCONFIRMED),
			NewBookingsByID(makeBooking(repUUID(32))),
			http.StatusOK,
			http.StatusOK,
			api.BookingStatusCONFIRMED,
		},

		{
			"posting a bookingEvent on existing booking (status CONFIRMED over CONFIRMED) fails with code 400",
			repUUID(34),
			makeCarpoolBookingEventWithStatus(repUUID(35), repUUID(34), api.BookingStatusCONFIRMED),
			NewBookingsByID(makeBookingWithStatus(repUUID(34), api.BookingStatusCONFIRMED)),
			http.StatusBadRequest,
			http.StatusOK,
			api.BookingStatusCONFIRMED,
		},

		{
			"posting a bookingEvent on existing booking (status CONFIRMED over CANCELLED) fails with code 400",
			repUUID(36),
			makeCarpoolBookingEventWithStatus(repUUID(37), repUUID(36), api.BookingStatusCONFIRMED),
			NewBookingsByID(makeBookingWithStatus(repUUID(36), api.BookingStatusCANCELLED)),
			http.StatusBadRequest,
			http.StatusOK,
			api.BookingStatusCANCELLED,
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			mockDB := db.NewMockDB()
			mockDB.Bookings = tc.existingBookings

			flagsPost := test.NewFlags()
			flagsPost.ExpectedResponseCode = tc.expectedPostStatusCode

			flagsGet := test.NewFlags()
			flagsGet.ExpectedResponseCode = tc.expectedGetStatusCode
			flagsGet.ExpectedBookingStatus = tc.expectedBookingStatus

			TestPostBookingEventsHelper(t, mockDB, *tc.carpoolBookingEvent, flagsPost)

			TestGetBookingsHelper(t, mockDB, tc.bookingID, flagsGet)
		})
	}
}

func TestPostMessage(t *testing.T) {
	var (
		bob    = makeUser("1", "bob")
		alice  = makeUser("2", "alice")
		carole = makeUser("3", "carole")
		david  = makeUser("4", "david")
		eve    = makeUser("5", "eve")
		fanny  = makeUser("6", "fanny")
	)

	testCases := []struct {
		name               string
		message            api.PostMessagesJSONBody
		existingUsers      []api.User
		expectedStatusCode int
	}{
		{
			"Posting message with both user known succeeds with code 201",
			makeMessage(alice, bob),
			[]api.User{alice, bob},
			http.StatusCreated,
		},

		{
			"Posting message with recipient unknown fails with code 404",
			makeMessage(carole, david),
			[]api.User{carole},
			http.StatusNotFound,
		},

		{
			"Posting message with sender unknown succeeds with code 201",
			makeMessage(eve, fanny),
			[]api.User{fanny},
			http.StatusCreated,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockDB := db.NewMockDB()
			mockDB.Users = tc.existingUsers

			flags := test.NewFlags()
			flags.ExpectedResponseCode = tc.expectedStatusCode

			TestPostMessagesHelper(t, mockDB, tc.message, flags)
		})
	}
}

func TestDefaultDriverJourneysValidity(t *testing.T) {
	params := requestAll(t, "driver")
	mockDB := db.NewMockDBWithDefaultData()
	testFun := test.TestGetDriverJourneysResponse

	flags := test.NewFlags()
	flags.ExpectNonEmpty = true

	testGetJourneysHelper(t, params, mockDB, testFun, flags)
}

func TestDefaultPassengerJourneysValidity(t *testing.T) {
	params := requestAll(t, "passenger")
	mockDB := db.NewMockDBWithDefaultData()
	testFun := test.TestGetPassengerJourneysResponse

	flags := test.NewFlags()
	flags.ExpectNonEmpty = true

	testGetJourneysHelper(t, params, mockDB, testFun, flags)
}

// Should be kept at the end as it relies on order of execution of tests.
func TestGeneration(t *testing.T) {
	if generateTestData {

		var b bytes.Buffer

		err := db.WriteData(generatedData, &b)
		util.PanicIf(err)

		bytes, err := io.ReadAll(&b)
		util.PanicIf(err)

		err = os.WriteFile(generatedTestDataFile, bytes, 0644)
		util.PanicIf(err)

		err = os.WriteFile(generatedTestCommandsFile, []byte(commandsFile.String()), 0644)
		util.PanicIf(err)
	}
}
