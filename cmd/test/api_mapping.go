package test

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// APIMapping is the mapping between endpoint and the associated test
// function that has been registered, if any
var APIMapping = map[Endpoint]ResponseTestFun{}

// Register associates a test function to a given function. If any
// TestFunction is already associated, it overwrites it.
func Register(f ResponseTestFun, e Endpoint) {
	APIMapping[e] = f
}

// Endpoint describes an Endpoint
type Endpoint struct {
	Method string
	Path   string
}

// String implements the Stringer interface for Endpoint type
func (e Endpoint) String() string {
	return e.Method + " " + e.Path
}

// SelectTestFuns returns the test functions related to a given request.
func SelectTestFuns(endpoint Endpoint) (ResponseTestFun, error) {
	testFun, ok := APIMapping[endpoint]
	if !ok {
		return nil, fmt.Errorf("request to an unknown endpoint: %s", endpoint)
	}
	return testFun, nil
}

// ExtractEndpoint extracts the endpoint from a request, given server
// information
func ExtractEndpoint(request *http.Request, server string) (Endpoint, error) {
	serverURL, err := url.Parse(server)
	if err != nil {
		return Endpoint{}, err
	}
	method := request.Method
	path := strings.TrimPrefix(request.URL.Path, serverURL.Path)
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return Endpoint{method, path}, nil
}
