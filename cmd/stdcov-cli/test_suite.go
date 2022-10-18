package main

import (
	"context"
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
func ExecuteTestSuite(client APIClient) Report {
	all := []AssertionResult{}
	for _, testFun := range TestSuite {
		all = append(all, testFun(client)...)
	}
	return Report{allAssertionResults: all}
}

/////////////////////////////////////////////////////////////

// A TestFun is a function that the API in a specific way (e.g. testing a
// single endpoint).
type TestFun func(APIClient) []AssertionResult

// TestGetStatus checks the `GET /status` endpoint
func TestGetStatus(Client APIClient) []AssertionResult {
	endpoint := Endpoint{"/status", http.MethodGet}
	a := NewAssertionAccu()
	a.endpoint = endpoint
	testGetStatus(Client, a)
	return a.GetAssertionResults()
}

// TestGetDriverJourneys checks the `GET /driver_journeys` endpoint
func TestGetDriverJourneys(Client APIClient) []AssertionResult {
	endpoint := Endpoint{"/driver_journeys", http.MethodGet}
	a := NewAssertionAccu()
	a.endpoint = endpoint
	testGetDriverJourneys(Client, a)
	return a.GetAssertionResults()
}

/////////////////////////////////////////////////////////////

type auxTestFun func(APIClient, AssertionAccumulator)

func testGetStatus(Client APIClient, a AssertionAccumulator) {
	response, clientErr := Client.GetStatus(context.Background())
	a.Run(
		Critic(assertAPICallSuccess{clientErr}),
		assertStatusCode{response, http.StatusOK},
	)
}

func testGetDriverJourneys(Client APIClient, a AssertionAccumulator) {
	// Test query parameters
	params := &client.GetDriverJourneysParams{}

	// Request
	request, _ := client.NewGetDriverJourneysRequest(Client.Server, params)

	// Get response
	response, clientErr := Client.GetDriverJourneys(context.Background(), params)

	a.Run(
		Critic(assertAPICallSuccess{clientErr}),
		assertStatusCode{response, http.StatusOK},
		assertHeaderContains{response, "Content-Type", "application/json"},
		assertDriverJourneysFormat{request, response},
	)
}
