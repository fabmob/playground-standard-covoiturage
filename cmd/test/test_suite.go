package test

import (
	"net/http"

	"gitlab.com/multi/stdcov-api-test/cmd/test/client"
)

// APIClient is a client to the API standard covoiturage
type APIClient = *client.Client

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
	tests []TestFun,
	flags Flags,
) *Report {
	all := []AssertionResult{}
	for _, testFun := range tests {
		all = append(all, testFun(client, request, flags)...)
	}
	return &Report{allAssertionResults: all}
}

/////////////////////////////////////////////////////////////

// A TestFun is a function that the API in a specific way (e.g. testing a
// single endpoint). Assumes that request is non-nil.
type TestFun func(APIClient, *http.Request, Flags) []AssertionResult

// wrapTest wraps an auxTestFun (that tests a response against a request) to a
// TestFun
func wrapTest(f auxTestFun, endpoint Endpoint) TestFun {
	return func(c APIClient, request *http.Request, flags Flags) []AssertionResult {
		a := NewAssertionAccu()
		a.endpoint = endpoint
		response, clientErr := c.Client.Do(request)
		if clientErr != nil {
			a.Run(assertAPICallSuccess{clientErr})
		} else {
			f(request, response, a, flags)
		}
		return a.GetAssertionResults()
	}
}

//////////////////////////////////////////////////////////////

type auxTestFun func(
	*http.Request,
	*http.Response,
	AssertionAccumulator,
	Flags,
)

func TestGetStatus(
	request *http.Request,
	response *http.Response,
	a AssertionAccumulator,
	flags Flags,
) {
	assertions := []Assertion{
		assertStatusCode{response, http.StatusOK},
	}
	a.Run(assertions...)
}

func TestGetDriverJourneys(
	request *http.Request,
	response *http.Response,
	a AssertionAccumulator,
	flags Flags,
) {

	response.Body = ReusableReadCloser(response.Body)

	var assertions []Assertion
	if flags.disallowEmpty {
		assertions = []Assertion{
			assertStatusCode{response, http.StatusOK},
			assertHeaderContains{response, "Content-Type", "application/json"},
			Critic(assertDriverJourneysNotEmpty{response}),
			Critic(assertDriverJourneysFormat{request, response}),
			assertDriverJourneysRadius{request, response, arrival}, assertDriverJourneysRadius{request, response, departure},
		}
	} else {
		assertions = []Assertion{
			assertStatusCode{response, http.StatusOK},
			assertHeaderContains{response, "Content-Type", "application/json"},
			Critic(assertDriverJourneysFormat{request, response}),
			assertDriverJourneysRadius{request, response, arrival}, assertDriverJourneysRadius{request, response, departure},
		}
	}

	a.Run(assertions...)
}