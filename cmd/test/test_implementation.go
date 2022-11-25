package test

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type testImplementation func(
	*http.Request,
	*http.Response,
	AssertionAccumulator,
	Flags,
)

//////////////////////////////////////////////////////////////
// Tag "Search"
//////////////////////////////////////////////////////////////

// TestGetDriverJourneys .. Assumes non empty response.
func testGetDriverJourneys(
	request *http.Request,
	response *http.Response,
	a AssertionAccumulator,
	flags Flags,
) {

	CriticAssertFormat(a, request, response)
	AssertStatusCode(a, response, flags.ExpectedResponseCode)
	AssertHeaderContains(a, response, echo.HeaderContentType, echo.MIMEApplicationJSON)

	if flags.ExpectNonEmpty {
		CriticAssertArrayNotEmpty(a, response)
	}

	AssertJourneysDepartureRadius(a, request, response)
	AssertJourneysArrivalRadius(a, request, response)
	AssertJourneysTimeDelta(a, request, response)
	AssertJourneysCount(a, request, response)
	AssertUniqueIDs(a, response)
	AssertOperatorFieldFormat(a, response)
}

func testGetPassengerJourneys(
	request *http.Request,
	response *http.Response,
	a AssertionAccumulator,
	flags Flags,
) {
	// Passenger journeys are very similar to driver journeys.
	testGetDriverJourneys(request, response, a, flags)
}

func testGetDriverRegularTrips(
	request *http.Request,
	response *http.Response,
	a AssertionAccumulator,
	flags Flags,
) {
	CriticAssertFormat(a, request, response)
	AssertStatusCode(a, response, flags.ExpectedResponseCode)
}

func testGetPassengerRegularTrips(
	request *http.Request,
	response *http.Response,
	a AssertionAccumulator,
	flags Flags,
) {
	CriticAssertFormat(a, request, response)
	AssertStatusCode(a, response, flags.ExpectedResponseCode)
}

//////////////////////////////////////////////////////////////
// Tag "Webhooks"
//////////////////////////////////////////////////////////////

func testPostBookingEvents(
	request *http.Request,
	response *http.Response,
	a AssertionAccumulator,
	flags Flags,
) {
	CriticAssertFormat(a, request, response)
	AssertStatusCode(a, response, flags.ExpectedResponseCode)
}

//////////////////////////////////////////////////////////////
// Tag "Interact"
//////////////////////////////////////////////////////////////

func testPostMessages(
	request *http.Request,
	response *http.Response,
	a AssertionAccumulator,
	flags Flags,
) {
	CriticAssertFormat(a, request, response)
	AssertStatusCode(a, response, flags.ExpectedResponseCode)
}

func testPostBookings(
	request *http.Request,
	response *http.Response,
	a AssertionAccumulator,
	flags Flags,
) {
	CriticAssertFormat(a, request, response)
	AssertStatusCode(a, response, flags.ExpectedResponseCode)
}

func testPatchBookings(
	request *http.Request,
	response *http.Response,
	a AssertionAccumulator,
	flags Flags,
) {
	CriticAssertFormat(a, request, response)
	AssertStatusCode(a, response, flags.ExpectedResponseCode)
}

// testGetBookings currently assumes that the request returns a 200 response.
func testGetBookings(
	request *http.Request,
	response *http.Response,
	a AssertionAccumulator,
	flags Flags,
) {

	CriticAssertFormat(a, request, response)
	AssertStatusCode(a, response, flags.ExpectedResponseCode)

	if flags.ExpectedBookingStatus != "" {
		AssertBookingStatus(a, response, string(flags.ExpectedBookingStatus))
	}
}

//////////////////////////////////////////////////////////////
// Tag "status"
//////////////////////////////////////////////////////////////

func testGetStatus(
	request *http.Request,
	response *http.Response,
	a AssertionAccumulator,
	flags Flags,
) {
	AssertStatusCodeOK(a, response)
}
