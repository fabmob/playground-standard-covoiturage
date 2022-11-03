package test

import "net/http"

type testImplementation func(
	*http.Request,
	*http.Response,
	AssertionAccumulator,
	Flags,
)

func testGetStatus(
	request *http.Request,
	response *http.Response,
	a AssertionAccumulator,
	flags Flags,
) {
	AssertStatusCodeOK(a, response)
}

// TestGetDriverJourneys .. Assumes non empty response.
func testGetDriverJourneys(
	request *http.Request,
	response *http.Response,
	a AssertionAccumulator,
	flags Flags,
) {

	AssertStatusCodeOK(a, response)
	AssertHeaderContains(a, response, "Content-Type", "application/json")
	if flags.DisallowEmpty {
		CriticAssertArrayNotEmpty(a, response)
	}
	CriticAssertFormat(a, request, response)
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
