package main

import (
	"fmt"
	"net/http"
)

// Endpoint describes an Endpoint
type Endpoint struct {
	method string
	path   string
}

// String implements the Stringer interface for Endpoint type
func (e Endpoint) String() string {
	return e.method + " " + e.path
}

// GetStatusEndpoint is the Endpoint of GET /status
var GetStatusEndpoint = Endpoint{http.MethodGet, "/status"}

// GetDriverJourneyEndpoint is the Endpoint of GET /driver_journeys
var GetDriverJourneyEndpoint = Endpoint{http.MethodGet, "/driver_journeys"}

var apiMapping = map[Endpoint][]TestFun{
	GetStatusEndpoint:        {wrapTest(testGetStatus, GetStatusEndpoint)},
	GetDriverJourneyEndpoint: {wrapTest(testGetDriverJourneys, GetDriverJourneyEndpoint)},
}

// SelectTestFuns returns the test functions related to a given request
func SelectTestFuns(request *http.Request) ([]TestFun, error) {
	method := request.Method
	path := request.URL.Path
	testFuns, ok := apiMapping[Endpoint{method, path}]
	if !ok {
		return nil, fmt.Errorf("request to an unknown endpoint: %s %s", method,
			path)
	}
	return testFuns, nil
}
