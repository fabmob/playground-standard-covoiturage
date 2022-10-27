package test

import (
	"net/http"

	"gitlab.com/multi/stdcov-api-test/cmd/api"
)

// APIClient is a client to the API standard covoiturage
type APIClient = *api.Client

// ExecuteTestSuite tests a client against all implemented tests
func ExecuteTestSuite(client APIClient, request *http.Request, flags Flags) (*Report, error) {
	selectedTestFuns, err := SelectTestFuns(request, client.Server)
	if err != nil {
		return nil, err
	}
	return executeTestFuns(client, request, selectedTestFuns, flags), nil
}

func executeTestFuns(
	client APIClient,
	request *http.Request,
	tests []TestRequestFun,
	flags Flags,
) *Report {
	all := []AssertionResult{}
	for _, testFun := range tests {
		all = append(all, testFun(client, request, flags)...)
	}
	return &Report{allAssertionResults: all}
}

/////////////////////////////////////////////////////////////

// A TestRequestFun runs all tests associated with a given Request, and return
// the correspending `AssertionResult`s
type TestRequestFun func(APIClient, *http.Request, Flags) []AssertionResult

// wrapTestResponseFun wraps an TestResponseFun to a TestRequestFun
func wrapTestResponseFun(f TestResponseFun, endpoint Endpoint) TestRequestFun {
	return func(c APIClient, request *http.Request, flags Flags) []AssertionResult {
		a := NewAssertionAccu()
		a.endpoint = endpoint
		response, clientErr := c.Client.Do(request)
		if clientErr != nil {
			a.Queue(assertAPICallSuccess{clientErr})
			a.ExecuteAll()
			return a.GetAssertionResults()
		}
		return f(request, response, a, flags)
	}
}

//////////////////////////////////////////////////////////////

type TestResponseFun func(
	*http.Request,
	*http.Response,
	AssertionAccumulator,
	Flags,
) []AssertionResult

var (
	TestGetStatusResponse         TestResponseFun = wrapAssertionsFun(testGetStatus)
	TestGetDriverJourneysResponse                 = wrapAssertionsFun(testGetDriverJourneys)
)

func wrapAssertionsFun(f assertionFun) TestResponseFun {
	return func(
		req *http.Request,
		resp *http.Response,
		a AssertionAccumulator,
		flags Flags,
	) []AssertionResult {
		f(req, resp, a, flags)
		a.ExecuteAll()
		return a.GetAssertionResults()
	}
}

//////////////////////////////////////////////////////////////

type assertionFun func(
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
		CriticAssertDriverJourneysNotEmpty(a, response)
	}
	CriticAssertDriverJourneysFormat(a, request, response)
	AssertDriverJourneysDepartureRadius(a, request, response)
	AssertDriverJourneysArrivalRadius(a, request, response)
	AssertDriverJourneysTimeDelta(a, request, response)
	AssertDriverJourneysCount(a, request, response)
	AssertUniqueIDs(a, response)
	AssertOperatorFieldFormat(a, response)
}
