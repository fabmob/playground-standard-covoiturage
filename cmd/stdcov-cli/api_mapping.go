package main

import (
	"fmt"
	"net/http"
)

// Endpoint describes an Endpoint
type Endpoint struct {
	Method string
	Path   string
}

// String implements the Stringer interface for Endpoint type
func (e Endpoint) String() string {
	return e.Method + " " + e.Path
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
func SelectTestFuns(request *http.Request, server string) ([]TestFun, error) {
	testFuns, ok := apiMapping[ExtractEndpoint(request, server)]
	if !ok {
		return nil, fmt.Errorf("request to an unknown endpoint. Method: %s, path: %s",
			request.Method,
			request.URL.Path)
	}
	return testFuns, nil
}

func ExtractEndpoint(request *http.Request, server string) Endpoint {
	method := request.Method
	path := request.URL.Path
	return Endpoint{method, path}
}
