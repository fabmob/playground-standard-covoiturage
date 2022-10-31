package test

import "net/http"

// registerAllTests associates every available test to a given endpoint
func registerAllTests() {
	Register(TestGetDriverJourneysResponse, GetDriverJourneyEndpoint)
	Register(TestGetStatusResponse, GetStatusEndpoint)
}

//////////////////////////////////////////////////////////////
// Exported tests
//////////////////////////////////////////////////////////////

var (
	// TestGetStatusResponse tests response of GET /status
	TestGetStatusResponse ResponseTestFun = wrapAssertionsFun(testGetStatus)
	// TestGetDriverJourneysResponse tests response of GET /driver_journeys
	TestGetDriverJourneysResponse = wrapAssertionsFun(testGetDriverJourneys)
)

//////////////////////////////////////////////////////////////
// Endpoints
//////////////////////////////////////////////////////////////

var (
	// GetStatusEndpoint is the Endpoint of GET /status
	GetStatusEndpoint = Endpoint{http.MethodGet, "/status"}
	// GetDriverJourneyEndpoint is the Endpoint of GET /driver_journeys
	GetDriverJourneyEndpoint = Endpoint{http.MethodGet, "/driver_journeys"}
)
