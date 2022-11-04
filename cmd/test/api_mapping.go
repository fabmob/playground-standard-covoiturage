package test

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// apiMapping is the mapping between endpoint and the associated test
// function that has been registered, if any
var apiMapping = map[Endpoint]ResponseTestFun{}

func GetAPIMapping() map[Endpoint]ResponseTestFun {
	if len(apiMapping) == 0 {
		initAPIMapping()
	}
	return apiMapping
}

// Register associates a test function to a given function. If any
// TestFunction is already associated, it overwrites it.
func Register(f ResponseTestFun, e Endpoint) {
	apiMapping[e] = f
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

// emptyRequest returns an empty *http.Request to the endpoint
func (e Endpoint) emptyRequest() *http.Request {
	request, _ := http.NewRequest(e.Method, e.Path, nil)
	return request
}

// SelectTestFuns returns the test functions related to a given request.
func SelectTestFuns(endpoint Endpoint) (ResponseTestFun, error) {
	testFun, ok := GetAPIMapping()[endpoint]
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

// GuessServer try to guess the server, and returns server and path in case of
// success.
func GuessServer(method, URL string) (string, error) {
	u, err := url.Parse(URL)
	if err != nil {
		return "", err
	}

	uWithoutQuery := u
	uWithoutQuery.RawQuery = ""
	uWithoutQuery.Fragment = ""

	for endpoint := range GetAPIMapping() {
		if endpoint.Method == method && strings.HasSuffix(uWithoutQuery.String(), endpoint.Path) {
			fmt.Println("x")
			server := strings.TrimSuffix(uWithoutQuery.String(), endpoint.Path)
			return server, nil
		}
	}

	return "", fmt.Errorf(
		"did not recognize the endpoint with method %s in %s",
		method,
		uWithoutQuery,
	)
}
