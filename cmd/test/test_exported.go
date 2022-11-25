package test

import "github.com/fabmob/playground-standard-covoiturage/cmd/endpoint"

func init() {
	// Tag "Search"

	Register(TestGetDriverJourneysResponse, endpoint.GetDriverJourneys)
	Register(TestGetPassengerJourneysResponse, endpoint.GetPassengerJourneys)
	Register(TestGetDriverRegularTripsResponse, endpoint.GetDriverRegularTrips)
	Register(TestGetPassengerRegularTripsResponse, endpoint.GetPassengerRegularTrips)

	// Tag "Webhooks"

	Register(TestPostBookingEventsResponse, endpoint.PostBookingEvents)

	// Tag "Interact"

	Register(TestPostMessagesResponse, endpoint.PostMessages)
	Register(TestPostBookingsResponse, endpoint.PostBookings)
	Register(TestPatchBookingsResponse, endpoint.PatchBookings)
	Register(TestGetBookingsResponse, endpoint.GetBookings)

	// Tag "status"

	Register(TestGetStatusResponse, endpoint.GetStatus)
}

//////////////////////////////////////////////////////////////
// Exported tests
//////////////////////////////////////////////////////////////

var (
	// Tag "Search"

	// TestGetDriverJourneysResponse tests response of GET /driver_journeys
	TestGetDriverJourneysResponse = wrapTestImplementation(testGetDriverJourneys)
	// TestGetPassengerJourneysResponse tests response of GET /passenger_journeys
	TestGetPassengerJourneysResponse = wrapTestImplementation(testGetPassengerJourneys)
	// TestGetDriverRegularTripsResponse tests response of GET /driver_regular_trips
	TestGetDriverRegularTripsResponse = wrapTestImplementation(testGetDriverRegularTrips)
	// TestGetPassengerRegularTripsResponse tests response of GET /passenger_regular_trips
	TestGetPassengerRegularTripsResponse = wrapTestImplementation(testGetPassengerRegularTrips)

	// Tag "Webhooks"

	// TestPostBookingsResponse tests response of POST /bookings/{booking_id}
	TestPostBookingEventsResponse = wrapTestImplementation(testPostBookingEvents)

	// Tag "Interact"

	// TestPostMessagesResponse tests response of POST /messages
	TestPostMessagesResponse = wrapTestImplementation(testPostMessages)
	// TestPostBookingsResponse tests response of POST /bookings/{booking_id}
	TestPostBookingsResponse = wrapTestImplementation(testPostBookings)
	// TestPatchBookingsResponse tests response of PATCH /bookings/{booking_id}
	TestPatchBookingsResponse = wrapTestImplementation(testPatchBookings)
	// TestGetBookingsResponse tests response of GET /bookings/{booking_id}
	TestGetBookingsResponse = wrapTestImplementation(testGetBookings)

	// Tag "status"

	// TestGetStatusResponse tests response of GET /status
	TestGetStatusResponse ResponseTestFun = wrapTestImplementation(testGetStatus)
)
