package test

import (
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/fabmob/playground-standard-covoiturage/cmd/test/endpoint"
)

// apiMapping stores api endpoint > test functions data
var apiMapping = map[endpoint.Info]ResponseTestFun{}

// GetAPIMapping returns the mapping between endpoint and the associated test
// function
func GetAPIMapping() map[endpoint.Info]ResponseTestFun {
	return apiMapping
}

// Register associates a test function to a given function. If any
// TestFunction is already associated, it overwrites it.
func Register(f ResponseTestFun, e endpoint.Info) {
	apiMapping[e] = f
}

// SelectTestFuns returns the test functions related to a given request.
func SelectTestFuns(endpoint endpoint.Info) (ResponseTestFun, error) {
	testFun, ok := GetAPIMapping()[endpoint]
	if !ok {
		return nil, fmt.Errorf("request to an unknown endpoint: %s", endpoint)
	}

	return testFun, nil
}

// SplitServerEndpoint try to guess the server, and returns server and path in case of
// success.
func SplitServerEndpoint(method, URL string) (string, endpoint.Info, error) {
	u, err := url.Parse(URL)
	if err != nil {
		return "", endpoint.Info{}, err
	}

	removeQuery(u)

	for endpoint := range GetAPIMapping() {

		suffix := knownEndpointSuffix(u.String(), endpoint)
		if endpoint.Method == method && suffix != "" {
			server := strings.TrimSuffix(u.String(), suffix)
			return server, endpoint, nil
		}
	}

	return "", endpoint.Info{}, fmt.Errorf(
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
func knownEndpointSuffix(url string, endpoint endpoint.Info) string {
	var param string
	if endpoint.HasPathParam {
		url, param = path.Split(url)
		url = ensureNoTrailingSlash(url)
		param = ensureLeadingSlash(param)
	}

	if strings.HasSuffix(url, endpoint.Path) {
		return endpoint.Path + param
	}

	return ""
}

func ensureNoTrailingSlash(s string) string {
	return strings.TrimSuffix(s, "/")
}

func ensureLeadingSlash(path string) string {
	if strings.HasPrefix(path, "/") {
		return path
	}

	return "/" + path
}
