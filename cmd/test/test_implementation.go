package test

import (
	"fmt"
	"io"
	"net/http"
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
	AssertStatusCode(a, response, flags.ExpectedStatusCode)
	AssertHeaderContains(a, response, "Content-Type", "application/json")

	if flags.DisallowEmpty {
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
	AssertStatusCode(a, response, flags.ExpectedStatusCode)
}

func testGetPassengerRegularTrips(
	request *http.Request,
	response *http.Response,
	a AssertionAccumulator,
	flags Flags,
) {
	CriticAssertFormat(a, request, response)
	AssertStatusCode(a, response, flags.ExpectedStatusCode)
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
	AssertStatusCode(a, response, flags.ExpectedStatusCode)
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
	AssertStatusCode(a, response, flags.ExpectedStatusCode)
}

func testPostBookings(
	request *http.Request,
	response *http.Response,
	a AssertionAccumulator,
	flags Flags,
) {

	fmt.Println("-----------------")
	fmt.Printf("%+v\n", response)
	body, _ := io.ReadAll(response.Body)
	fmt.Printf("Body:\n %+v\n", string(body))

	CriticAssertFormat(a, request, response)
	AssertStatusCode(a, response, flags.ExpectedStatusCode)
}

func testPatchBookings(
	request *http.Request,
	response *http.Response,
	a AssertionAccumulator,
	flags Flags,
) {
	CriticAssertFormat(a, request, response)
	AssertStatusCode(a, response, flags.ExpectedStatusCode)
}

// testGetBookings currently assumes that the request returns a 200 response.
func testGetBookings(
	request *http.Request,
	response *http.Response,
	a AssertionAccumulator,
	flags Flags,
) {

	CriticAssertFormat(a, request, response)
	AssertStatusCode(a, response, flags.ExpectedStatusCode)

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
