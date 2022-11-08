package test

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
)

// apiMapping stores api endpoint > test functions data
var apiMapping = map[Endpoint]ResponseTestFun{}

// GetAPIMapping returns the mapping between endpoint and the associated test
// function
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
	Method       string
	Path         string
	HasPathParam bool
}

func NewEndpoint(method, path string) Endpoint {
	return Endpoint{method, path, false}
}

func NewEndpointWithParam(method, path string) Endpoint {
	return Endpoint{method, path, true}
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
	path = ensureLeadingSlash(path)
	firstPathSegment := firstPathSegment(path)

	return NewEndpoint(method, firstPathSegment), nil
}

func ensureLeadingSlash(path string) string {
	if strings.HasPrefix(path, "/") {
		return path
	}

	return "/" + path
}

// firstPathSegment assumes without checking that path has a leading slash
func firstPathSegment(path string) string {
	return "/" + strings.Split(path, "/")[1]
}

// SplitServerEndpoint try to guess the server, and returns server and path in case of
// success.
func SplitServerEndpoint(method, URL string) (string, Endpoint, error) {
	u, err := url.Parse(URL)
	if err != nil {
		return "", Endpoint{}, err
	}

	removeQuery(u)

	for endpoint := range GetAPIMapping() {

		suffix := knownEndpointSuffix(u.String(), endpoint)
		if endpoint.Method == method && suffix != "" {
			server := strings.TrimSuffix(u.String(), suffix)
			return server, endpoint, nil
		}
	}

	return "", Endpoint{}, fmt.Errorf(
		"did not recognize the endpoint with method %s in %s",
		method,
		u,
	)
}

func removeQuery(u *url.URL) {
	u.RawQuery = ""
	u.Fragment = ""
}

// knownEndpointSuffix returns the complete endpoint suffix (including path
// parameter) if the endpoint path is recognized, an empty string
// otherwise.
func knownEndpointSuffix(url string, endpoint Endpoint) string {
	var param string
	if endpoint.HasPathParam {
		url, param = path.Split(url)
		url = removeTrailingSlash(url)
		param = ensureLeadingSlash(param)
	}

	if strings.HasSuffix(url, endpoint.Path) {
		return endpoint.Path + param
	}

	return ""
}

func removeTrailingSlash(s string) string {
	return strings.TrimSuffix(s, "/")
}
