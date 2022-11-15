package cmd

import (
	"net/http"
	"testing"

	"github.com/fabmob/playground-standard-covoiturage/cmd/test"
)

func TestPatchBookingsCmd(t *testing.T) {

	var (
		server    = "https://localhost:9999"
		bookingID = "9999"
		status    = "CONFIRMED"
		message   = "test message"
		expected  = expectedData{
			method:            http.MethodPatch,
			url:               "https://localhost:9999/bookings/9999",
			defaultStatusCode: http.StatusOK,
			body:              nil,
		}
	)

	mockRunner := test.NewMockRunner()
	err := patchBookingsRun(mockRunner, server, bookingID, status, message)
	panicIf(err)

	// Test Assertions
	expected.testArgs(t, mockRunner)

	testQueryParam(t, mockRunner.Query, "status", status)
	testQueryParam(t, mockRunner.Query, "message", message)
}

func TestPostMessagesCmd(t *testing.T) {

	var (
		server   = "https://localhost:9999"
		body     = []byte("body")
		expected = expectedData{
			method:            http.MethodPost,
			url:               "https://localhost:9999/messages",
			defaultStatusCode: http.StatusCreated,
			body:              body,
		}
	)

	mockRunner := test.NewMockRunner()
	err := getMessagesRun(mockRunner, server, body)
	panicIf(err)

	// Test Assertions
	expected.testArgs(t, mockRunner)
}

func TestPostBookingEventsCmd(t *testing.T) {

	var (
		server   = "https://localhost:9999"
		body     = []byte("body")
		expected = expectedData{
			method:            http.MethodPost,
			url:               "https://localhost:9999/booking_events",
			defaultStatusCode: http.StatusOK,
			body:              body,
		}
	)

	mockRunner := test.NewMockRunner()
	err := postBookingEventsRun(mockRunner, server, body)
	panicIf(err)

	expected.testArgs(t, mockRunner)
}

func TestGetDriverRegularTripsCmd(t *testing.T) {
	var (
		server             = "https://localhost:9999"
		departureLat       = "0"
		departureLng       = "1"
		arrivalLat         = "2"
		arrivalLng         = "3"
		departureTimeOfDay = "4"
		timeDelta          = "5"
		departureRadius    = "6"
		arrivalRadius      = "7"
		count              = "8"
		expected           = expectedData{
			method:            http.MethodGet,
			url:               "https://localhost:9999/driver_regular_trips",
			defaultStatusCode: http.StatusOK,
			body:              nil,
		}
	)

	mockRunner := test.NewMockRunner()
	err := getDriverRegularTripsRun(
		mockRunner,
		server,
		departureLat, departureLng, arrivalLat, arrivalLng, departureTimeOfDay,
		timeDelta, departureRadius, arrivalRadius, count,
	)
	panicIf(err)

	expected.testArgs(t, mockRunner)
}

func TestGetPassengerRegularTripsCmd(t *testing.T) {
	var (
		server             = "https://localhost:9999"
		departureLat       = "0"
		departureLng       = "1"
		arrivalLat         = "2"
		arrivalLng         = "3"
		departureTimeOfDay = "4"
		departureWeekdays  = "4"
		timeDelta          = "5"
		departureRadius    = "6"
		arrivalRadius      = "7"
		count              = "8"
		expected           = expectedData{
			method:            http.MethodGet,
			url:               "https://localhost:9999/passenger_regular_trips",
			defaultStatusCode: http.StatusOK,
			body:              nil,
		}
	)

	mockRunner := test.NewMockRunner()
	err := getPassengerRegularTripsRun(
		mockRunner,
		server,
		departureLat, departureLng, arrivalLat, arrivalLng, departureTimeOfDay,
		timeDelta, departureRadius, arrivalRadius, count,
	)
	panicIf(err)

	expected.testArgs(t, mockRunner)

	testQueryParam(t, mockRunner.Query, "departureTimeOfDay", departureTimeOfDay)
	testQueryParam(t, mockRunner.Query, "departureWeekdays", departureWeekdays)
	/* testQueryParam(t, mockRunner.Query, "minDepartureDate", minDepartureDate) */
	/* testQueryParam(t, mockRunner.Query, "maxDepartureDate", maxDepartureDate) */
}

func testStringArg(t *testing.T, got, expected, argumentName string) {
	t.Helper()
	if expected != got {
		t.Logf("Unexpected %s in command.", argumentName)
		t.Logf("Expected %s", expected)
		t.Logf("Got %s", got)
		t.Fail()
	}
}

func testIntArg(t *testing.T, got, expected int, argumentName string) {
	t.Helper()
	if expected != got {
		t.Logf("Unexpected %s in command.", argumentName)
		t.Logf("Expected %d", expected)
		t.Logf("Got %d", got)
		t.Fail()
	}
}

func testQueryParam(t *testing.T, query test.Query, param, value string) {
	gotValue, ok := query.Params[param]
	if !ok {
		t.Errorf("Missing query parameter %s", param)
	} else {
		testStringArg(t, gotValue, value, param)
	}
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

type expectedData struct {
	method            string
	url               string
	defaultStatusCode int
	body              []byte
}

func (expected expectedData) testArgs(t *testing.T, mockRunner *test.MockRunner) {
	testStringArg(t, mockRunner.Method, expected.method, "method")

	testStringArg(t, mockRunner.URL, expected.url, "URL")

	testIntArg(t, mockRunner.Flags.ExpectedStatusCode,
		expected.defaultStatusCode, "status code")

	nilBodyExpected := expected.body == nil
	nilBodyProvided := mockRunner.Body == nil

	if nilBodyExpected && !nilBodyProvided {
		t.Error("Body provided while none expected")
	} else if !nilBodyExpected && nilBodyProvided {
		t.Error("Required body is missing")
	} else if !nilBodyExpected && !nilBodyProvided {
		testStringArg(t, string(mockRunner.Body), string(expected.body), "body")
	}
}
