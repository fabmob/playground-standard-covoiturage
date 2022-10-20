package main

import (
	"net/http"

	"gitlab.com/multi/stdcov-api-test/cmd/stdcov-cli/client"
)

// APIClient is a client to the API standard covoiturage
type APIClient = *client.Client

// TestSuite lists all test functions that are executed when the API is tested
var TestSuite = []TestFun{
	TestGetStatus,
	TestGetDriverJourneys,
}

// ExecuteTestSuite tests a client against all implemented tests
func ExecuteTestSuite(client APIClient, request *http.Request) Report {
	all := []AssertionResult{}
	for _, testFun := range TestSuite {
		all = append(all, testFun(client, request)...)
	}
	return Report{allAssertionResults: all}
}

/////////////////////////////////////////////////////////////

// A TestFun is a function that the API in a specific way (e.g. testing a
// single endpoint). Assumes that request is non-nil.
type TestFun func(APIClient, *http.Request) []AssertionResult

// TestGetStatus tests the GET /status endpoint
var TestGetStatus = wrapTest(testGetStatus)

// TestGetDriverJourneys tests the GET /driver_journeys endpoint
var TestGetDriverJourneys = wrapTest(testGetDriverJourneys)

// wrapTest wraps an auxTestFun (that tests a response against a request) to a
// TestFun
func wrapTest(f auxTestFun) TestFun {
	return func(c APIClient, request *http.Request) []AssertionResult {
		endpoint := Endpoint{request.URL.Path, request.Method}
		a := NewAssertionAccu()
		a.endpoint = endpoint
		response, clientErr := c.Client.Do(request)
		if clientErr != nil {
			a.Run(assertAPICallSuccess{clientErr})
		} else {
			f(request, response, a)
		}
		return a.GetAssertionResults()
	}
}

//////////////////////////////////////////////////////////////

type auxTestFun func(*http.Request, *http.Response, AssertionAccumulator)

func testGetStatus(
	request *http.Request,
	response *http.Response,
	a AssertionAccumulator,
) {

	a.Run(
		assertStatusCode{response, http.StatusOK},
	)
}

func testGetDriverJourneys(
	request *http.Request,
	response *http.Response,
	a AssertionAccumulator,
) {

	a.Run(
		assertStatusCode{response, http.StatusOK},
		assertHeaderContains{response, "Content-Type", "application/json"},
		assertDriverJourneysFormat{request, response},
	)
}
