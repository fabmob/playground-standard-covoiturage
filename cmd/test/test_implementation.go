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
	response.Body = ReusableReadCloser(response.Body)

	AssertStatusCodeOK(a, response)
	AssertHeaderContains(a, response, "Content-Type", "application/json")
	if flags.DisallowEmpty {
		CriticAssertArrayNotEmpty(a, response)
	}
	CriticAssertFormat(a, request, response)
	AssertJourneysDepartureRadius(a, request, response)
	AssertJourneysArrivalRadius(a, request, response)
	AssertDriverJourneysTimeDelta(a, request, response)
	AssertDriverJourneysCount(a, request, response)
	AssertUniqueIDs(a, response)
	AssertOperatorFieldFormat(a, response)
}

func testGetPassengerJourneys(
	request *http.Request,
	response *http.Response,
	a AssertionAccumulator,
	flags Flags,
) {
	response.Body = ReusableReadCloser(response.Body)

	AssertStatusCodeOK(a, response)
	AssertHeaderContains(a, response, "Content-Type", "application/json")
	if flags.DisallowEmpty {
		CriticAssertArrayNotEmpty(a, response)
	}
	CriticAssertFormat(a, request, response)
	AssertJourneysDepartureRadius(a, request, response)
	AssertJourneysArrivalRadius(a, request, response)
	AssertDriverJourneysTimeDelta(a, request, response)
	AssertDriverJourneysCount(a, request, response)
	AssertUniqueIDs(a, response)
	AssertOperatorFieldFormat(a, response)
}
