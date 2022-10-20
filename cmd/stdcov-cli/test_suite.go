package main

import (
	"context"
	"net/http"
	"strings"

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
func ExecuteTestSuite(client APIClient) Report {
	all := []AssertionResult{}
	request, _ := http.NewRequest("GET", "/", strings.NewReader(""))
	for _, testFun := range TestSuite {
		all = append(all, testFun(client, request)...)
	}
	return Report{allAssertionResults: all}
}

/////////////////////////////////////////////////////////////

// A TestFun is a function that the API in a specific way (e.g. testing a
// single endpoint).
type TestFun func(APIClient, *http.Request) []AssertionResult

// TestGetStatus checks the `GET /status` endpoint
func TestGetStatus(Client APIClient, request *http.Request) []AssertionResult {
	endpoint := Endpoint{"/status", http.MethodGet}
	a := NewAssertionAccu()
	a.endpoint = endpoint
	testGetStatus(Client, request, a)
	return a.GetAssertionResults()
}

// TestGetDriverJourneys checks the `GET /driver_journeys` endpoint
func TestGetDriverJourneys(Client APIClient, request *http.Request) []AssertionResult {
	endpoint := Endpoint{"/driver_journeys", http.MethodGet}
	a := NewAssertionAccu()
	a.endpoint = endpoint
	testGetDriverJourneys(Client, request, a)
	return a.GetAssertionResults()
}

//////////////////////////////////////////////////////////////

type auxTestFun func(APIClient, *http.Request, AssertionAccumulator)

func testGetStatus(Client APIClient, request *http.Request, a AssertionAccumulator) {
	response, clientErr := Client.GetStatus(context.Background())
	a.Run(
		Critic(assertAPICallSuccess{clientErr}),
		assertStatusCode{response, http.StatusOK},
	)
}

func testGetDriverJourneys(c APIClient, request *http.Request, a AssertionAccumulator) {
	request, _ = http.NewRequest(http.MethodGet, "/driver_journeys", strings.NewReader(""))

	// Get response
	response, clientErr := c.Client.Do(request)

	a.Run(
		Critic(assertAPICallSuccess{clientErr}),
		assertStatusCode{response, http.StatusOK},
		assertHeaderContains{response, "Content-Type", "application/json"},
		assertDriverJourneysFormat{request, response},
	)
}
