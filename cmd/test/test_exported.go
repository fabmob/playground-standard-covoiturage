package test

import "net/http"

// initAPIMapping associates every available test to a given endpoint
func initAPIMapping() {
	Register(TestGetStatusResponse, GetStatusEndpoint)
	Register(TestGetDriverJourneysResponse, GetDriverJourneysEndpoint)
	Register(TestGetPassengerJourneysResponse, GetPassengerJourneysEndpoint)
	Register(TestGetBookingsResponse, GetBookingsEndpoint)
}

//////////////////////////////////////////////////////////////
// Exported tests
//////////////////////////////////////////////////////////////

var (
	// TestGetStatusResponse tests response of GET /status
	TestGetStatusResponse ResponseTestFun = wrapAssertionsFun(testGetStatus)
	// TestGetDriverJourneysResponse tests response of GET /driver_journeys
	TestGetDriverJourneysResponse = wrapAssertionsFun(testGetDriverJourneys)
	// TestGetPassengerJourneysResponse tests response of GET /passenger_journeys
	TestGetPassengerJourneysResponse = wrapAssertionsFun(testGetPassengerJourneys)
	// TestGetBookingsResponse tests response of GET /bookings/{booking_id}
	TestGetBookingsResponse = wrapAssertionsFun(testGetBookings)
)

//////////////////////////////////////////////////////////////
// Endpoints
//////////////////////////////////////////////////////////////

var (
	// GetStatusEndpoint is the Endpoint of GET /status
	GetStatusEndpoint = NewEndpoint(http.MethodGet, "/status")
	// GetDriverJourneysEndpoint is the Endpoint of GET /driver_journeys
	GetDriverJourneysEndpoint = NewEndpoint(http.MethodGet, "/driver_journeys")
	// GetPassengerJourneysEndpoint is the Endpoint of GET /passenger_journeys
	GetPassengerJourneysEndpoint = NewEndpoint(http.MethodGet, "/passenger_journeys")
	// GetBookingsEndpoint is the Endpoint of GET /passenger_journeys
	GetBookingsEndpoint = NewEndpointWithParam(http.MethodGet, "/bookings")
)
