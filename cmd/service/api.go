package service

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"gitlab.com/multi/stdcov-api-test/cmd/service/server"
)

// StdCovServerImpl implements server.ServerInterface
type StdCovServerImpl struct {
	mockDB mockDB
}

func NewDefaultServer() (*StdCovServerImpl, error) {
	server := StdCovServerImpl{NewMockDB()}
	err := server.mockDB.PopulateDBWithDefault()
	return &server, err
}

// PostBookingEvents sends booking information of a user connected with a third-party provider back to the provider.
// (POST /booking_events)
func (*StdCovServerImpl) PostBookingEvents(ctx echo.Context) error {
	// Implement me
	return nil
}

// PostBookings creates a punctual outward Booking requet.
// (POST /bookings)
func (*StdCovServerImpl) PostBookings(ctx echo.Context) error {
	// Implement me
	return nil
}

// GetBookings retrieves an existing Booking request.
// (GET /bookings/{bookingId})
func (*StdCovServerImpl) GetBookings(ctx echo.Context, bookingID server.BookingId) error {
	// Implement me
	return nil
}

// PatchBookings updates status of an existing Booking request.
// (PATCH /bookings/{bookingId})
func (*StdCovServerImpl) PatchBookings(ctx echo.Context, bookingID server.BookingId,
	params server.PatchBookingsParams) error {
	// Implement me
	return nil
}

// GetDriverJourneys searches for matching punctual planned outward driver journeys.
// (GET /driver_journeys)
func (s *StdCovServerImpl) GetDriverJourneys(
	ctx echo.Context,
	params server.GetDriverJourneysParams,
) error {
	return ctx.JSON(http.StatusOK, []server.DriverJourney{})
}

// GetDriverRegularTrips searches for matching regular driver trip.
// (GET /driver_regular_trips)
func (*StdCovServerImpl) GetDriverRegularTrips(
	ctx echo.Context,
	params server.GetDriverRegularTripsParams,
) error {
	// Implement me
	return nil
}

// PostConnections sends a mesage to the owner of a retrieved journey.
// (POST /messages)
func (*StdCovServerImpl) PostConnections(ctx echo.Context) error {
	// Implement me
	return nil
}

// GetPassengerJourneys searches for matching punctual planned outward pasenger journeys.
// (GET /passenger_journeys)
func (*StdCovServerImpl) GetPassengerJourneys(
	ctx echo.Context,
	params server.GetPassengerJourneysParams,
) error {
	// Implement me
	return nil
}

// GetPassengerRegularTrips searches for matching pasenger regular trips.
// (GET /passenger_regular_trips)
func (*StdCovServerImpl) GetPassengerRegularTrips(
	ctx echo.Context,
	params server.GetPassengerRegularTripsParams,
) error {
	// Implement me
	return nil
}

// GetStatus gives health status of the webservice.
// (GET /status)
func (*StdCovServerImpl) GetStatus(ctx echo.Context) error {
	return ctx.NoContent(http.StatusOK)
}
