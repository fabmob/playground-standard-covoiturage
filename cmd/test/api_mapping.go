package test

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
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
	GetStatusEndpoint:        {wrapTest(TestGetStatus, GetStatusEndpoint)},
	GetDriverJourneyEndpoint: {wrapTest(TestGetDriverJourneys, GetDriverJourneyEndpoint)},
}

// SelectTestFuns returns the test functions related to a given request
func SelectTestFuns(request *http.Request, server string) ([]TestFun, error) {
	endpoint, err := ExtractEndpoint(request, server)
	if err != nil {
		return nil, err
	}
	testFuns, ok := apiMapping[*endpoint]
	if !ok {
		return nil, fmt.Errorf("request to an unknown endpoint. Method: %s, path: %s",
			request.Method,
			request.URL.Path)
	}
	return testFuns, nil
}

// ExtractEndpoint extracts the endpoint from a request, given server
// information
func ExtractEndpoint(request *http.Request, server string) (*Endpoint, error) {
	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}
	method := request.Method
	path := strings.TrimPrefix(request.URL.Path, serverURL.Path)
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return &Endpoint{method, path}, nil
}