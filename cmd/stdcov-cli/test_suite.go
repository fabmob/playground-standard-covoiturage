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
	a := NewDefaultAsserter()
	a.endpoint = endpoint
	testGetStatus(Client, a)
	return a.GetAssertionResults()
}

// TestGetDriverJourneys checks the `GET /driver_journeys` endpoint
func TestGetDriverJourneys(Client APIClient) []AssertionResult {
	endpoint := Endpoint{"/driver_journeys", http.MethodGet}
	a := NewDefaultAsserter()
	a.endpoint = endpoint
	testGetDriverJourneys(Client, a)
	return a.GetAssertionResults()
}

/////////////////////////////////////////////////////////////

type auxTestFun func(APIClient, AssertionAccumulator)

func testGetStatus(Client APIClient, a AssertionAccumulator) {
	response, err := Client.GetStatus(context.Background())

	AssertAPICallSuccess(a, err)
	if a.LastAssertionHasError() {
		return
	}

	AssertStatusCodeOK(a, response)
}

func testGetDriverJourneys(Client APIClient, a AssertionAccumulator) {
	// Test query parameters
	params := &client.GetDriverJourneysParams{}

	// Request
	request, err := client.NewGetDriverJourneysRequest(Client.Server, params)
	AssertAPICallSuccess(a, err)

	// Get response
	response, err := Client.GetDriverJourneys(context.Background(), params)
	AssertAPICallSuccess(a, err)
	if a.LastAssertionHasError() {
		return
	}

	AssertStatusCodeOK(a, response)
	if a.LastAssertionHasError() {
		return
	}
	AssertHeaderContains(a, response, "Content-Type", "application/json")
	AssertDriverJourneysFormat(a, request, response)
}
