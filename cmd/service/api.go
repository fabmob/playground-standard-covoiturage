package service

import (
	"errors"
	"net/http"

	"github.com/fabmob/playground-standard-covoiturage/cmd/api"
	"github.com/fabmob/playground-standard-covoiturage/cmd/service/db"
	"github.com/labstack/echo/v4"
)

// StdCovServerImpl implements server.ServerInterface
type StdCovServerImpl struct {
	db db.DB
}

func NewServer() *StdCovServerImpl {
	server := StdCovServerImpl{db.NewMockDB()}
	return &server
}

func NewServerWithDB(mockDB db.DB) *StdCovServerImpl {
	server := StdCovServerImpl{mockDB}
	return &server
}

// NewDefaultServer returns a server, and populates the associated DB with
// default data
func NewDefaultServer() *StdCovServerImpl {
	server := StdCovServerImpl{db.NewMockDBWithDefaultData()}
	return &server
}

// PostBookingEvents sends booking information of a user connected with a third-party provider back to the provider.
// (POST /booking_events)
func (s *StdCovServerImpl) PostBookingEvents(ctx echo.Context) error {
	var newEvent api.CarpoolBookingEvent

	bodyUnmarshallingErr := ctx.Bind(&newEvent)
	if bodyUnmarshallingErr != nil {
		return ctx.JSON(http.StatusBadRequest, errorBody(bodyUnmarshallingErr))
	}

	var newBooking api.Booking

	if driverCarpoolBooking, err := newEvent.Data.AsDriverCarpoolBooking(); err == nil {
		newBooking = *driverCarpoolBooking.ToBooking()
	} else if passengerCarpoolBooking, err := newEvent.Data.AsPassengerCarpoolBooking(); err == nil {
		newBooking = *passengerCarpoolBooking.ToBooking()
	} else {
		return ctx.JSON(
			http.StatusBadRequest,
			errorBody(errors.New("unmarshaling error")),
		)
	}

	// Try to add booking
	alreadyExistsErr := s.db.AddBooking(newBooking)

	// If booking exists, try to update status
	if alreadyExistsErr != nil {
		err := UpdateBookingStatus(s.db, newBooking.Id, newBooking.Status)

		if err != nil {
			var missing db.MissingBookingErr

			switch {
			case errors.As(err, &missing):
				// should not happen
				return ctx.NoContent(http.StatusInternalServerError)

			default:
				return ctx.JSON(http.StatusBadRequest, errorBody(err))
			}
		}
	}

	return ctx.NoContent(http.StatusOK)
}

// PostBookings creates a punctual outward Booking request.
// (POST /bookings)
func (s *StdCovServerImpl) PostBookings(ctx echo.Context) error {
	var newBooking api.Booking

	// Unmarshal request body into newBooking
	bodyUnmarshallingErr := ctx.Bind(&newBooking)
	if bodyUnmarshallingErr != nil {
		return ctx.JSON(http.StatusBadRequest, errorBody(bodyUnmarshallingErr))
	}

	alreadyExistsErr := s.db.AddBooking(newBooking)
	if alreadyExistsErr != nil {
		return ctx.JSON(http.StatusBadRequest, errorBody(alreadyExistsErr))
	}

	return ctx.JSON(http.StatusCreated, newBooking)
}

type Error struct {
	Error string `json:"error"`
}

// GetBookings retrieves an existing Booking request.
// (GET /bookings/{bookingId})
func (s *StdCovServerImpl) GetBookings(ctx echo.Context, bookingID api.BookingId) error {

	booking, missingErr := s.db.GetBooking(bookingID)

	if missingErr != nil {
		return ctx.JSON(http.StatusNotFound, errorBody(missingErr))
	}

	return ctx.JSON(http.StatusOK, booking)
}

// PatchBookings updates status of an existing Booking request.
// (PATCH /bookings/{bookingId})
func (s *StdCovServerImpl) PatchBookings(ctx echo.Context, bookingID api.BookingId,
	params api.PatchBookingsParams) error {

	err := UpdateBookingStatus(s.db, bookingID, params.Status)

	if err != nil {
		var missing db.MissingBookingErr
		var statusAlreadySet StatusAlreadySetErr

		switch {
		case errors.As(err, &missing):
			return ctx.JSON(http.StatusNotFound, errorBody(err))

		case errors.As(err, &statusAlreadySet):
			return ctx.JSON(http.StatusConflict, errorBody(err))

		default:
			return ctx.JSON(http.StatusBadRequest, errorBody(err))
		}
	}

	return ctx.NoContent(http.StatusOK)
}

// GetDriverJourneys searches for matching punctual planned outward driver journeys.
// (GET /driver_journeys)
func (s *StdCovServerImpl) GetDriverJourneys(
	ctx echo.Context,
	params api.GetDriverJourneysParams,
) error {
	response := []api.DriverJourney{}

	for _, dj := range s.db.GetDriverJourneys() {
		if keepJourney(&params, dj.Trip, dj.JourneySchedule) {
			response = append(response, dj)
		}
	}

	if params.Count != nil {
		response = keepNFirst(response, *params.Count)
	}

	return ctx.JSON(http.StatusOK, response)
}

// GetDriverRegularTrips searches for matching regular driver trip.
// (GET /driver_regular_trips)
func (s *StdCovServerImpl) GetDriverRegularTrips(
	ctx echo.Context,
	params api.GetDriverRegularTripsParams,
) error {
	response := []api.DriverRegularTrip{}

	for _, drt := range s.db.GetDriverRegularTrips() {
		if !keepTrip(&params, drt.Trip) {
			continue
		}

		schedulesOK, err := anyScheduleOK(drt.Schedules, &params)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, errorBody(err))
		}
		if schedulesOK {
			response = append(response, drt)
		}
	}

	if params.Count != nil {
		response = keepNFirst(response, *params.Count)
	}

	return ctx.JSON(http.StatusOK, response)
}

// PostMessages sends a mesage to the owner of a retrieved journey.
// (POST /messages)
func (s *StdCovServerImpl) PostMessages(ctx echo.Context) error {
	users := s.db.GetUsers()

	var message api.PostMessagesJSONBody

	bodyUnmarshallingErr := ctx.Bind(&message)
	if bodyUnmarshallingErr != nil {
		return ctx.JSON(http.StatusBadRequest, errorBody(bodyUnmarshallingErr))
	}

	if !userExists(message.To, users) {
		return ctx.JSON(http.StatusNotFound, errorBody(errors.New("missing_user")))
	}

	return ctx.NoContent(http.StatusCreated)
}

// GetPassengerJourneys searches for matching punctual planned outward pasenger journeys.
// (GET /passenger_journeys)
func (s *StdCovServerImpl) GetPassengerJourneys(
	ctx echo.Context,
	params api.GetPassengerJourneysParams,
) error {
	response := []api.PassengerJourney{}

	for _, pj := range s.db.GetPassengerJourneys() {
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
func (s *StdCovServerImpl) GetPassengerRegularTrips(
	ctx echo.Context,
	params api.GetPassengerRegularTripsParams,
) error {
	response := []api.PassengerRegularTrip{}

	for _, prt := range s.db.GetPassengerRegularTrips() {
		if !keepTrip(&params, prt.Trip) {
			continue
		}

		schedulesOK, err := anyScheduleOK(prt.Schedules, &params)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, errorBody(err))
		}
		if schedulesOK {
			response = append(response, prt)
		}
	}

	if params.Count != nil {
		response = keepNFirst(response, *params.Count)
	}

	return ctx.JSON(http.StatusOK, response)
}

// GetStatus gives health status of the webservice.
// (GET /status)
func (*StdCovServerImpl) GetStatus(ctx echo.Context) error {
	return ctx.NoContent(http.StatusOK)
}
