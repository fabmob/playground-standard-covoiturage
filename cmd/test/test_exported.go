package test

import "net/http"

// initAPIMapping associates every available test to a given endpoint
func initAPIMapping() {
	// Tag "Search"

	Register(TestGetDriverJourneysResponse, GetDriverJourneysEndpoint)
	Register(TestGetPassengerJourneysResponse, GetPassengerJourneysEndpoint)
	Register(TestGetDriverRegularTripsResponse, GetDriverRegularTripsEndpoint)
	Register(TestGetPassengerRegularTripsResponse, GetPassengerRegularTripsEndpoint)

	// Tag "Webhooks"

	Register(TestPostBookingEventsResponse, PostBookingEventsEndpoint)

	// Tag "Interact"

	Register(TestPostMessagesResponse, PostMessagesEndpoint)
	Register(TestPostBookingsResponse, PostBookingsEndpoint)
	Register(TestPatchBookingsResponse, PatchBookingsEndpoint)
	Register(TestGetBookingsResponse, GetBookingsEndpoint)

	// Tag "status"

	Register(TestGetStatusResponse, GetStatusEndpoint)
}

//////////////////////////////////////////////////////////////
// Exported tests
//////////////////////////////////////////////////////////////

var (
	// Tag "Search"

	// TestGetDriverJourneysResponse tests response of GET /driver_journeys
	TestGetDriverJourneysResponse = wrapAssertionsFun(testGetDriverJourneys)
	// TestGetPassengerJourneysResponse tests response of GET /passenger_journeys
	TestGetPassengerJourneysResponse = wrapAssertionsFun(testGetPassengerJourneys)
	// TestGetDriverRegularTripsResponse tests response of GET /driver_regular_trips
	TestGetDriverRegularTripsResponse = wrapAssertionsFun(testGetDriverRegularTrips)
	// TestGetPassengerRegularTripsResponse tests response of GET /passenger_regular_trips
	TestGetPassengerRegularTripsResponse = wrapAssertionsFun(testGetPassengerRegularTrips)

	// Tag "Webhooks"

	// TestPostBookingsResponse tests response of POST /bookings/{booking_id}
	TestPostBookingEventsResponse = wrapAssertionsFun(testPostBookingEvents)

	// Tag "Interact"

	// TestPostMessagesResponse tests response of POST /messages
	TestPostMessagesResponse = wrapAssertionsFun(testPostMessages)
	// TestPostBookingsResponse tests response of POST /bookings/{booking_id}
	TestPostBookingsResponse = wrapAssertionsFun(testPostBookings)
	// TestPatchBookingsResponse tests response of PATCH /bookings/{booking_id}
	TestPatchBookingsResponse = wrapAssertionsFun(testPatchBookings)
	// TestGetBookingsResponse tests response of GET /bookings/{booking_id}
	TestGetBookingsResponse = wrapAssertionsFun(testGetBookings)

	// Tag "status"

	// TestGetStatusResponse tests response of GET /status
	TestGetStatusResponse ResponseTestFun = wrapAssertionsFun(testGetStatus)
)

//////////////////////////////////////////////////////////////
// Endpoints
//////////////////////////////////////////////////////////////

var (
	// Tag "Search"

	// GetDriverJourneysEndpoint is the Endpoint of GET /driver_journeys
	GetDriverJourneysEndpoint = NewEndpoint(http.MethodGet, "/driver_journeys")
	// GetPassengerJourneysEndpoint is the Endpoint of GET /passenger_journeys
	GetPassengerJourneysEndpoint = NewEndpoint(http.MethodGet, "/passenger_journeys")
	// GetDriverRegularTripsEndpoint is the Endpoint of GET /driver_regularTrips
	GetDriverRegularTripsEndpoint = NewEndpoint(http.MethodGet, "/driver_regular_trips")
	// GetPassengerRegularTripsEndpoint is the Endpoint of GET /passenger_regularTrips
	GetPassengerRegularTripsEndpoint = NewEndpoint(http.MethodGet, "/passenger_regular_trips")

	// Tag "Webhooks"

	// PostBookingEventsEndpoint is the Endpoint of POST /passenger_journeys
	PostBookingEventsEndpoint = NewEndpoint(http.MethodPost, "/booking_events")

	// Tag "Interact"

	// PostMessagesEndpoint is the Endpoint of POST /messages
	PostMessagesEndpoint = NewEndpoint(http.MethodPost, "/messages")
	// PostBookingsEndpoint is the Endpoint of POST /passenger_journeys
	PostBookingsEndpoint = NewEndpoint(http.MethodPost, "/bookings")
	// PatchBookingsEndpoint is the Endpoint of PATCH /passenger_journeys
	PatchBookingsEndpoint = NewEndpointWithParam(http.MethodPatch, "/bookings")
	// GetBookingsEndpoint is the Endpoint of GET /passenger_journeys
	GetBookingsEndpoint = NewEndpointWithParam(http.MethodGet, "/bookings")

	// Tag "status"

	// GetStatusEndpoint is the Endpoint of GET /status
	GetStatusEndpoint = NewEndpoint(http.MethodGet, "/status")
)
