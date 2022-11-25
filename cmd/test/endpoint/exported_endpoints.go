package endpoint

import "net/http"

var (
	// Tag "Search"

	// GetDriverJourneys is the of GET /driver_journeys
	GetDriverJourneys = New(http.MethodGet, "/driver_journeys")
	// GetPassengerJourneys is the of GET /passenger_journeys
	GetPassengerJourneys = New(http.MethodGet, "/passenger_journeys")
	// GetDriverRegularTrips is the of GET /driver_regularTrips
	GetDriverRegularTrips = New(http.MethodGet, "/driver_regular_trips")
	// GetPassengerRegularTrips is the of GET /passenger_regularTrips
	GetPassengerRegularTrips = New(http.MethodGet, "/passenger_regular_trips")

	// Tag "Webhooks"

	// PostBookingEvents is the of POST /passenger_journeys
	PostBookingEvents = New(http.MethodPost, "/booking_events")

	// Tag "Interact"

	// PostMessages is the of POST /messages
	PostMessages = New(http.MethodPost, "/messages")
	// PostBookings is the of POST /passenger_journeys
	PostBookings = New(http.MethodPost, "/bookings")
	// PatchBookings is the of PATCH /passenger_journeys
	PatchBookings = NewWithParam(http.MethodPatch, "/bookings")
	// GetBookings is the of GET /passenger_journeys
	GetBookings = NewWithParam(http.MethodGet, "/bookings")

	// Tag "status"

	// GetStatus is the of GET /status
	GetStatus = New(http.MethodGet, "/status")
)
