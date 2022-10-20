package main

import "net/http"

// GetStatusEndpoint is the Endpoint of GET /status
var GetStatusEndpoint = Endpoint{http.MethodGet, "/status"}

// GetDriverJourneyEndpoint is the Endpoint of GET /driver_journeys
var GetDriverJourneyEndpoint = Endpoint{http.MethodGet, "/driver_journeys"}

var apiMapping = map[Endpoint][]TestFun{
	GetStatusEndpoint:        {wrapTest(testGetStatus, GetStatusEndpoint)},
	GetDriverJourneyEndpoint: {wrapTest(testGetDriverJourneys, GetDriverJourneyEndpoint)},
}
