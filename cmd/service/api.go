package service

import (
	"math"
	"net/http"

	"github.com/fabmob/playground-standard-covoiturage/cmd/api"
	"github.com/fabmob/playground-standard-covoiturage/cmd/util"
	"github.com/labstack/echo/v4"
)

// StdCovServerImpl implements server.ServerInterface
type StdCovServerImpl struct {
	mockDB *MockDB
}

func NewServer() *StdCovServerImpl {
	server := StdCovServerImpl{NewMockDB()}
	return &server
}

func NewServerWithDB(mockDB *MockDB) *StdCovServerImpl {
	server := StdCovServerImpl{mockDB}
	return &server
}

// NewDefaultServer returns a server, and populates the associated DB with
// default data
func NewDefaultServer() *StdCovServerImpl {
	server := StdCovServerImpl{NewMockDBWithDefaultData()}
	return &server
}

// PostBookingEvents sends booking information of a user connected with a third-party provider back to the provider.
// (POST /booking_events)
func (*StdCovServerImpl) PostBookingEvents(ctx echo.Context) error {
	// Implement me
	return nil
}

// PostBookings creates a punctual outward Booking request.
// (POST /bookings)
func (s *StdCovServerImpl) PostBookings(ctx echo.Context) error {
	var newBooking *api.Booking
	err := ctx.Bind(&newBooking)
	if err != nil {
		errorStr := err.Error()
		return ctx.JSON(
			http.StatusBadRequest,
			api.BadRequest{Error: &errorStr},
		)
	}

	if _, idExists := s.mockDB.Bookings[newBooking.Id]; idExists {
		errorStr := "booking ID already exists"
		return ctx.JSON(http.StatusBadRequest, api.BadRequest{Error: &errorStr})
	}

	s.mockDB.Bookings[newBooking.Id] = newBooking

	return ctx.JSON(http.StatusCreated, newBooking)
}

type Error struct {
	Error string `json:"error"`
}

// GetBookings retrieves an existing Booking request.
// (GET /bookings/{bookingId})
func (s *StdCovServerImpl) GetBookings(ctx echo.Context, bookingID api.BookingId) error {

	for _, booking := range s.mockDB.Bookings {
		if booking.Id == bookingID {
			return ctx.JSON(http.StatusOK, booking)
		}
	}

	errReturn := Error{"missing_booking"}
	return ctx.JSON(http.StatusNotFound, errReturn)
}

// PatchBookings updates status of an existing Booking request.
// (PATCH /bookings/{bookingId})
func (s *StdCovServerImpl) PatchBookings(ctx echo.Context, bookingID api.BookingId,
	params api.PatchBookingsParams) error {

	booking, ok := s.mockDB.Bookings[bookingID]
	if !ok {
		errorStr := "missing_booking"
		return ctx.JSON(http.StatusNotFound, api.BadRequest{Error: &errorStr})
	}

	booking.Status = params.Status

	return nil
}

// GetDriverJourneys searches for matching punctual planned outward driver journeys.
// (GET /driver_journeys)
func (s *StdCovServerImpl) GetDriverJourneys(
	ctx echo.Context,
	params api.GetDriverJourneysParams,
) error {
	response := []api.DriverJourney{}

	for _, dj := range s.mockDB.DriverJourneys {
		if keepJourney(&params, dj.Trip, dj.JourneySchedule) {
			response = append(response, dj)
		}
	}

	if params.Count != nil {
		response = keepNFirst(response, *params.Count)
	}

	return ctx.JSON(http.StatusOK, response)
}

func keepJourney(params api.GetJourneysParams, trip api.Trip, schedule api.JourneySchedule) bool {
	coordsRequestDeparture := util.Coord{
		Lat: float64(params.GetDepartureLat()),
		Lon: float64(params.GetDepartureLng()),
	}
	coordsResponseDeparture := util.Coord{
		Lat: trip.PassengerPickupLat,
		Lon: trip.PassengerPickupLng,
	}
	departureRadiusOK := util.Distance(coordsRequestDeparture, coordsResponseDeparture) <=
		params.GetDepartureRadius()

	coordsRequestArrival := util.Coord{
		Lat: float64(params.GetArrivalLat()),
		Lon: float64(params.GetArrivalLng()),
	}
	coordsResponseArrival := util.Coord{
		Lat: trip.PassengerDropLat,
		Lon: trip.PassengerDropLng,
	}
	arrivalRadiusOK := util.Distance(coordsRequestArrival, coordsResponseArrival) <=
		params.GetArrivalRadius()

	timeDeltaOK :=
		math.Abs(float64(schedule.PassengerPickupDate)-float64(params.GetDepartureDate())) <
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
func (s *StdCovServerImpl) GetPassengerJourneys(
	ctx echo.Context,
	params api.GetPassengerJourneysParams,
) error {
	response := []api.PassengerJourney{}

	for _, pj := range s.mockDB.PassengerJourneys {
		if keepJourney(&params, pj.Trip, pj.JourneySchedule) {
			response = append(response, pj)
		}
	}

	if params.Count != nil {
		response = keepNFirst(response, *params.Count)
	}

	return ctx.JSON(http.StatusOK, response)
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
