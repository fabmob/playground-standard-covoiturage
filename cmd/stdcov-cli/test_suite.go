package main

import (
	"net/http"

	"gitlab.com/multi/stdcov-api-test/cmd/stdcov-cli/client"
)

// APIClient is a client to the API standard covoiturage
type APIClient = *client.Client

// ExecuteTestSuite tests a client against all implemented tests
func ExecuteTestSuite(client APIClient, request *http.Request) (*Report, error) {
	selectedTestFuns, err := SelectTestFuns(request, client.Server)
	if err != nil {
		return nil, err
	}
	return executeTestFuns(client, request, selectedTestFuns), nil
}

func executeTestFuns(client APIClient, request *http.Request, tests []TestFun) *Report {
	all := []AssertionResult{}
	for _, testFun := range tests {
		all = append(all, testFun(client, request)...)
	}
	return &Report{allAssertionResults: all}
}

/////////////////////////////////////////////////////////////

// A TestFun is a function that the API in a specific way (e.g. testing a
// single endpoint). Assumes that request is non-nil.
type TestFun func(APIClient, *http.Request) []AssertionResult

// wrapTest wraps an auxTestFun (that tests a response against a request) to a
// TestFun
func wrapTest(f auxTestFun, endpoint Endpoint) TestFun {
	return func(c APIClient, request *http.Request) []AssertionResult {
		a := NewAssertionAccu()
		a.endpoint = endpoint
		response, clientErr := c.Client.Do(request)
		response.Body = ReusableReadCloser(response.Body)
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
		Critic(assertDriverJourneysFormat{request, response}),
		assertDriverJourneysRadius{request, response, arrival},
		assertDriverJourneysRadius{request, response, departure},
	)
}
