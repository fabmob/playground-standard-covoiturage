package test

import (
	"fmt"

	"github.com/fabmob/playground-standard-covoiturage/cmd/endpoint"
)

// apiMapping stores api endpoint > test functions data
var apiMapping = map[endpoint.Info]ResponseTestFun{}

// GetAPIMapping returns the mapping between endpoint and the associated test
// function
func GetAPIMapping() map[endpoint.Info]ResponseTestFun {
	return apiMapping
}

// Register associates a test function to a given endpoint. If any
// TestFunction is already associated, it overwrites it.
func Register(f ResponseTestFun, e endpoint.Info) {
	apiMapping[e] = f
}

// SelectTestFuns returns the test function related to a given request.
func SelectTestFuns(endpoint endpoint.Info) (ResponseTestFun, error) {
	testFun, ok := GetAPIMapping()[endpoint]
	if !ok {
		return nil, fmt.Errorf("request to an unknown endpoint: %s", endpoint)
	}

	return testFun, nil
}
