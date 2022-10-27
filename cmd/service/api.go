package service

import (
	"math"
	"net/http"

	"github.com/labstack/echo/v4"
	"gitlab.com/multi/stdcov-api-test/cmd/api"
	"gitlab.com/multi/stdcov-api-test/cmd/util"
)

// StdCovServerImpl implements server.ServerInterface
type StdCovServerImpl struct {
	mockDB MockDB
}

// NewDefaultServer returns a server, and populates the associated DB with
// default data
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
func (*StdCovServerImpl) GetBookings(ctx echo.Context, bookingID api.BookingId) error {
	// Implement me
	return nil
}

// PatchBookings updates status of an existing Booking request.
// (PATCH /bookings/{bookingId})
func (*StdCovServerImpl) PatchBookings(ctx echo.Context, bookingID api.BookingId,
	params api.PatchBookingsParams) error {
	// Implement me
	return nil
}

// GetDriverJourneys searches for matching punctual planned outward driver journeys.
// (GET /driver_journeys)
func (s *StdCovServerImpl) GetDriverJourneys(
	ctx echo.Context,
	params api.GetDriverJourneysParams,
) error {
	response := []api.DriverJourney{}
	for _, dj := range s.mockDB.driverJourneys {
		if keepDriverJourney(params, dj) {
			response = append(response, dj)
		}
	}
	if params.Count != nil {
		response = response[0:*params.Count]
	}
	return ctx.JSON(http.StatusOK, response)
}

func keepDriverJourney(params api.GetDriverJourneysParams, dj api.DriverJourney) bool {
	coordsRequestDeparture := util.Coord{
		Lat: float64(params.DepartureLat),
		Lon: float64(params.DepartureLng),
	}
	coordsResponseDeparture := util.Coord{
		Lat: dj.PassengerPickupLat,
		Lon: dj.PassengerPickupLng,
	}
	departureRadiusOK := util.Distance(coordsRequestDeparture, coordsResponseDeparture) <=
		params.GetDepartureRadius()

	coordsRequestArrival := util.Coord{
		Lat: float64(params.ArrivalLat),
		Lon: float64(params.ArrivalLng),
	}
	coordsResponseArrival := util.Coord{
		Lat: dj.PassengerDropLat,
		Lon: dj.PassengerDropLng,
	}
	arrivalRadiusOK := util.Distance(coordsRequestArrival, coordsResponseArrival) <=
		params.GetArrivalRadius()

	timeDeltaOK :=
		math.Abs(float64(dj.PassengerPickupDate)-float64(params.DepartureDate)) <
			float64(params.GetTimeDelta())
	return departureRadiusOK && arrivalRadiusOK && timeDeltaOK
}

// GetDriverRegularTrips searches for matching regular driver trip.
// (GET /driver_regular_trips)
func (*StdCovServerImpl) GetDriverRegularTrips(
	ctx echo.Context,
	params api.GetDriverRegularTripsParams,
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
	params api.GetPassengerJourneysParams,
) error {
	// Implement me
	return nil
}

// GetPassengerRegularTrips searches for matching pasenger regular trips.
// (GET /passenger_regular_trips)
func (*StdCovServerImpl) GetPassengerRegularTrips(
	ctx echo.Context,
	params api.GetPassengerRegularTripsParams,
) error {
	// Implement me
	return nil
}

// GetStatus gives health status of the webservice.
// (GET /status)
func (*StdCovServerImpl) GetStatus(ctx echo.Context) error {
	return ctx.NoContent(http.StatusOK)
}
