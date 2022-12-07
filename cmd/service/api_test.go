package service

import (
	"bytes"
	"flag"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/fabmob/playground-standard-covoiturage/cmd/api"
	"github.com/fabmob/playground-standard-covoiturage/cmd/service/db"
	"github.com/fabmob/playground-standard-covoiturage/cmd/test"
	"github.com/fabmob/playground-standard-covoiturage/cmd/util"
	"github.com/google/uuid"
)

func init() {
	// test flags do not need to be parsed explicitely, as it is already done in
	// normal `go test` operation
	flag.BoolVar(&generateTestData, "generate", false, "Should test data be regenerated")
}

var (
	coordsIgnore = util.Coord{Lat: 0, Lon: 0}
	coordsRef    = util.Coord{Lat: 46.1604531, Lon: -1.2219607} // reference
	coords900m   = util.Coord{Lat: 46.1613442, Lon: -1.2103736} // at ~900m from reference
	coords1100m  = util.Coord{Lat: 46.1613679, Lon: -1.2086563} // at ~1100m from reference
	coords2100m  = util.Coord{Lat: 46.1649225, Lon: -1.1954497} // at ~2100m from reference
)

// tripTestCases are common to GET /driver_journeys, GET /passenger_journeys,
// GET /driver_regular_trips and GET /passenger_regular_trips
var tripTestCases = []tripTestCase{
	{
		"No data",
		&api.GetDriverJourneysParams{},
		[]api.Trip{},
		false,
	},

	{
		"Departure radius 1",
		makeParamsWithDepartureRadius(coordsRef, 1, "driver"),
		[]api.Trip{
			makeTripAtCoords(coords900m, coordsIgnore),
			makeTripAtCoords(coords1100m, coordsIgnore),
		},
		true,
	},

	{
		"Departure radius 2",
		makeParamsWithDepartureRadius(coordsRef, 2, "driver"),
		[]api.Trip{
			makeTripAtCoords(coords900m, coordsIgnore),
			makeTripAtCoords(coords2100m, coordsIgnore),
		},
		true,
	},

	{
		"Departure radius 3",
		makeParamsWithDepartureRadius(coordsRef, 1, "driver"),
		[]api.Trip{
			makeTripAtCoords(coords1100m, coordsIgnore),
		},
		false,
	},

	{
		"Departure radius 3",
		makeParamsWithDepartureRadius(coordsRef, 1, "driver"),
		[]api.Trip{
			makeTripAtCoords(coords900m, coordsIgnore),
		},
		true,
	},

	{
		"Arrival radius 1",
		makeParamsWithArrivalRadius(coordsRef, 1, "driver"),
		[]api.Trip{
			makeTripAtCoords(coordsIgnore, coords900m),
			makeTripAtCoords(coordsIgnore, coords1100m),
		},
		true,
	},

	{
		"Arrival radius 2",
		makeParamsWithArrivalRadius(coordsRef, 2, "driver"),
		[]api.Trip{
			makeTripAtCoords(coordsIgnore, coords2100m),
			makeTripAtCoords(coordsIgnore, coords900m),
		},
		true,
	},

	{
		"Arrival radius 3",
		makeParamsWithArrivalRadius(coordsRef, 1, "driver"),
		[]api.Trip{
			makeTripAtCoords(coordsIgnore, coords1100m),
		},
		false,
	},

	{
		"Arrival radius 4",
		makeParamsWithArrivalRadius(coordsRef, 1, "driver"),
		[]api.Trip{
			makeTripAtCoords(coordsIgnore, coords900m),
		},
		true,
	},

	{
		"Count 1",
		makeParamsWithCount(1, "driver"),
		makeNTrips(1),
		true,
	},

	{
		"Count 2",
		makeParamsWithCount(0, "driver"),
		makeNTrips(1),
		false,
	},

	{
		"Count 3",
		makeParamsWithCount(2, "driver"),
		makeNTrips(4),
		true,
	},

	{
		"Count 4 - count > n driver journeys",
		makeParamsWithCount(1, "driver"),
		makeNTrips(0),
		false,
	},
}

// journeyScheduleTestCases are test cases common to GET /driver_journeys and
// GET /passenger_journeys
var journeyScheduleTestCases = []journeyScheduleTestCase{
	{
		"TimeDelta 1",
		makeParamsWithTimeDelta(10,
			"driver"),
		[]api.JourneySchedule{
			makeJourneyScheduleAtDate(5),
		},
		true,
	},

	{
		"TimeDelta 2",
		makeParamsWithTimeDelta(10,
			"driver"),
		[]api.JourneySchedule{
			makeJourneyScheduleAtDate(15),
		},
		false,
	},

	{
		"TimeDelta 3",
		makeParamsWithTimeDelta(20,
			"driver"),
		[]api.JourneySchedule{
			makeJourneyScheduleAtDate(25),
			makeJourneyScheduleAtDate(15),
		},
		true,
	},
}

func TestDriverJourneys(t *testing.T) {
	testCases := []driverJourneysTestCase{}

	for _, tc := range tripTestCases {
		testCases = append(testCases, tc.promoteToDriverJourneysTestCase(t))
	}
	for _, tc := range journeyScheduleTestCases {
		testCases = append(testCases, tc.promoteToDriverJourneysTestCase(t))
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			// If data is generated, then the test data and the requests date
			// properties are shifted, so that there are no two tests falling the
			// same week. The aim is to isolate the tests.
			if generateTestData {
				shiftToNextWeek()

				for i := range tc.testData {
					setJourneyDatesForGeneration(&tc.testData[i].JourneySchedule)
				}

				setParamDatesForGeneration(tc.testParams)
			}

			mockDB := db.NewMockDB()
			mockDB.DriverJourneys = tc.testData

			flags := test.NewFlags()
			flags.ExpectNonEmpty = tc.expectNonEmptyResult

			TestGetDriverJourneysHelper(
				t,
				mockDB,
				tc.testParams.(*api.GetDriverJourneysParams),
				flags,
			)
		})
	}
}

func TestPassengerJourneys(t *testing.T) {
	testCases := []passengerJourneysTestCase{}

	for _, tc := range tripTestCases {
		testCases = append(testCases, tc.promoteToPassengerJourneysTestCase(t))
	}
	for _, tc := range journeyScheduleTestCases {
		testCases = append(testCases, tc.promoteToPassengerJourneysTestCase(t))
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {

			// If data is generated, then the test data and the requests date
			// properties are shifted, so that there are no two tests falling the
			// same week. The aim is to isolate the tests.
			if generateTestData {
				shiftToNextWeek()

				for i := range tc.testData {
					setJourneyDatesForGeneration(&tc.testData[i].JourneySchedule)
				}

				setParamDatesForGeneration(tc.testParams)
			}

			mockDB := db.NewMockDB()
			mockDB.PassengerJourneys = tc.testData

			flags := test.NewFlags()
			flags.ExpectNonEmpty = tc.expectNonEmptyResult

			TestGetPassengerJourneysHelper(
				t,
				mockDB,
				tc.testParams.(*api.GetPassengerJourneysParams),
				flags,
			)
		})
	}
}

func TestGetDriverRegularTrips(t *testing.T) {
	testCases := []driverRegularTripsTestCase{}

	for _, tc := range tripTestCases {
		testCases = append(testCases, tc.promoteToDriverRegularTripsTestCase(t))
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			// If data is generated, then the test data and the requests date
			// properties are shifted, so that there are no two tests falling the
			// same week. The aim is to isolate the tests.
			if generateTestData {
				shiftToNextWeek()

				for i := range tc.testData {
					if tc.testData[i].Schedules != nil {
						schedules := *tc.testData[i].Schedules
						for _, schedule := range schedules {
							if schedule.JourneySchedules != nil {
								jschedules := *schedule.JourneySchedules
								for _, jschedule := range jschedules {
									setJourneyDatesForGeneration(&jschedule)
								}
							}
						}
					}
				}

				setParamDatesForGeneration(tc.testParams)
			}

			mockDB := db.NewMockDB()
			mockDB.DriverRegularTrips = tc.testData

			flags := test.NewFlags()
			flags.ExpectNonEmpty = tc.expectNonEmptyResult

			TestGetDriverRegularTripsHelper(
				t,
				mockDB,
				tc.testParams.(*api.GetDriverRegularTripsParams),
				flags,
			)
		})
	}
}

func TestGetPassengerRegularTrips(t *testing.T) {

	testCases := []struct {
		name                 string
		testParams           *api.GetPassengerRegularTripsParams
		testData             []api.PassengerRegularTrip
		expectNonEmptyResult bool
	}{
		{
			"No data",
			&api.GetPassengerRegularTripsParams{},
			[]api.PassengerRegularTrip{},
			false,
		},

		{
			"Valid regular trip",
			&api.GetPassengerRegularTripsParams{},
			[]api.PassengerRegularTrip{
				api.NewPassengerRegularTrip(),
			},
			true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			// If data is generated, then the test data and the requests date
			// properties are shifted, so that there are no two tests falling the
			// same week. The aim is to isolate the tests.
			if generateTestData {
				/* shiftToNextWeek() */

				/* for i := range tc.testData { */
				/* 	setJourneyDatesForGeneration(&tc.testData[i].JourneySchedule) */
				/* } */

				/* setParamDatesForGeneration(tc.testParams) */
			}

			mockDB := db.NewMockDB()
			mockDB.PassengerRegularTrips = tc.testData

			flags := test.NewFlags()
			flags.ExpectNonEmpty = tc.expectNonEmptyResult

			TestGetPassengerRegularTripsHelper(
				t,
				mockDB,
				tc.testParams,
				flags,
			)
		})
	}
}

func TestGetBookings(t *testing.T) {
	testCases := []struct {
		name               string
		bookings           db.BookingsByID
		queryBookingID     uuid.UUID
		expectNonEmpty     bool
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
			flags.ExpectNonEmpty = tc.expectNonEmpty
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

// Should be kept after tests that requires generation as it relies on order of execution of tests.
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
	generateTestData = false
}

// after this, no test is generated

func TestDefaultDriverJourneysValidity(t *testing.T) {
	params := requestAll(t, "driver")
	mockDB := db.NewMockDBWithDefaultData()

	flags := test.NewFlags()
	flags.ExpectNonEmpty = true

	TestGetDriverJourneysHelper(
		t,
		mockDB,
		params.(*api.GetDriverJourneysParams),
		flags,
	)
}

func TestDefaultPassengerJourneysValidity(t *testing.T) {
	params := requestAll(t, "passenger")
	mockDB := db.NewMockDBWithDefaultData()

	flags := test.NewFlags()
	flags.ExpectNonEmpty = true

	TestGetPassengerJourneysHelper(
		t,
		mockDB,
		params.(*api.GetPassengerJourneysParams),
		flags,
	)
}
