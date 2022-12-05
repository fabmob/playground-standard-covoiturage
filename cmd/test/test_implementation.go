package test

import (
	"net/http"

	"github.com/fabmob/playground-standard-covoiturage/cmd/test/assert"
	"github.com/labstack/echo/v4"
)

type testImplementation func(
	*http.Request,
	*http.Response,
	assert.Accumulator,
	Flags,
)

//////////////////////////////////////////////////////////////
// Tag "Search"
//////////////////////////////////////////////////////////////

// TestGetDriverJourneys .. Assumes non empty response.
func testGetDriverJourneys(
	request *http.Request,
	response *http.Response,
	a assert.Accumulator,
	flags Flags,
) {

	assert.CriticFormat(a, request, response)
	assert.StatusCode(a, response, flags.ExpectedResponseCode)
	assert.HeaderContains(a, response, echo.HeaderContentType, echo.MIMEApplicationJSON)

	if flags.ExpectNonEmpty {
		assert.CriticArrayNotEmpty(a, response)
	}

	assert.JourneysDepartureRadius(a, request, response)
	assert.JourneysArrivalRadius(a, request, response)
	assert.JourneysTimeDelta(a, request, response)
	assert.JourneysCount(a, request, response)
	assert.UniqueIDs(a, response)
	assert.OperatorFieldFormat(a, response)
}

func testGetPassengerJourneys(
	request *http.Request,
	response *http.Response,
	a assert.Accumulator,
	flags Flags,
) {
	// Passenger journeys are very similar to driver journeys.
	testGetDriverJourneys(request, response, a, flags)
}

func testGetDriverRegularTrips(
	request *http.Request,
	response *http.Response,
	a assert.Accumulator,
	flags Flags,
) {
	assert.CriticFormat(a, request, response)
	assert.StatusCode(a, response, flags.ExpectedResponseCode)

	if flags.ExpectNonEmpty {
		assert.CriticArrayNotEmpty(a, response)
	}

	assert.JourneysDepartureRadius(a, request, response)
	assert.JourneysArrivalRadius(a, request, response)
	assert.JourneysCount(a, request, response)
	assert.OperatorFieldFormat(a, response)
}

func testGetPassengerRegularTrips(
	request *http.Request,
	response *http.Response,
	a assert.Accumulator,
	flags Flags,
) {
	assert.CriticFormat(a, request, response)
	assert.StatusCode(a, response, flags.ExpectedResponseCode)

	if flags.ExpectNonEmpty {
		assert.CriticArrayNotEmpty(a, response)
	}
}

//////////////////////////////////////////////////////////////
// Tag "Webhooks"
//////////////////////////////////////////////////////////////

func testPostBookingEvents(
	request *http.Request,
	response *http.Response,
	a assert.Accumulator,
	flags Flags,
) {
	assert.CriticFormat(a, request, response)
	assert.StatusCode(a, response, flags.ExpectedResponseCode)
}

//////////////////////////////////////////////////////////////
// Tag "Interact"
//////////////////////////////////////////////////////////////

func testPostMessages(
	request *http.Request,
	response *http.Response,
	a assert.Accumulator,
	flags Flags,
) {
	assert.CriticFormat(a, request, response)
	assert.StatusCode(a, response, flags.ExpectedResponseCode)
}

func testPostBookings(
	request *http.Request,
	response *http.Response,
	a assert.Accumulator,
	flags Flags,
) {
	assert.CriticFormat(a, request, response)
	assert.StatusCode(a, response, flags.ExpectedResponseCode)
}

func testPatchBookings(
	request *http.Request,
	response *http.Response,
	a assert.Accumulator,
	flags Flags,
) {
	assert.CriticFormat(a, request, response)
	assert.StatusCode(a, response, flags.ExpectedResponseCode)
}

// testGetBookings currently assumes that the request returns a 200 response.
func testGetBookings(
	request *http.Request,
	response *http.Response,
	a assert.Accumulator,
	flags Flags,
) {

	assert.CriticFormat(a, request, response)
	assert.StatusCode(a, response, flags.ExpectedResponseCode)

	if flags.ExpectedBookingStatus != "" {
		assert.BookingStatus(a, response, string(flags.ExpectedBookingStatus))
	}
}

//////////////////////////////////////////////////////////////
// Tag "status"
//////////////////////////////////////////////////////////////

func testGetStatus(
	request *http.Request,
	response *http.Response,
	a assert.Accumulator,
	flags Flags,
) {
	assert.StatusCodeOK(a, response)
}
