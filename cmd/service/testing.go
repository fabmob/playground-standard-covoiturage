package service

import (
	"errors"
	"math/rand"
	"net/http"
	"net/http/httptest"

	"github.com/fabmob/playground-standard-covoiturage/cmd/api"
	"github.com/fabmob/playground-standard-covoiturage/cmd/util"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const fakeServer = "http://localhost:1323"

func setupTestServer(
	db *MockDB,
	request *http.Request,
) (*StdCovServerImpl, echo.Context, *httptest.ResponseRecorder) {

	var (
		e       = echo.New()
		rec     = httptest.NewRecorder()
		ctx     = e.NewContext(request, rec)
		handler = NewServerWithDB(db)
	)

	return handler, ctx, rec
}

func makeNDriverJourneys(n int) []api.DriverJourney {
	driverJourneys := make([]api.DriverJourney, 0, n)

	for i := 0; i < n; i++ {
		driverJourneys = append(driverJourneys, api.NewDriverJourney())
	}

	return driverJourneys
}

func makeNPassengerJourneys(n int) []api.PassengerJourney {
	passengerJourneys := make([]api.PassengerJourney, 0, n)

	for i := 0; i < n; i++ {
		passengerJourneys = append(passengerJourneys, api.NewPassengerJourney())
	}

	return passengerJourneys
}

func makeDriverJourneyAtCoords(coordPickup, coordDrop util.Coord) api.DriverJourney {
	dj := api.NewDriverJourney()
	updateTripCoords(&dj.Trip, coordPickup, coordDrop)

	return dj
}

func makePassengerJourneyAtCoords(coordPickup, coordDrop util.Coord) api.PassengerJourney {
	pj := api.NewPassengerJourney()
	updateTripCoords(&pj.Trip, coordPickup, coordDrop)

	return pj
}

func updateTripCoords(t *api.Trip, coordPickup, coordDrop util.Coord) {
	t.PassengerPickupLat = coordPickup.Lat
	t.PassengerPickupLng = coordPickup.Lon
	t.PassengerDropLat = coordDrop.Lat
	t.PassengerDropLng = coordDrop.Lon
}

func makeDriverJourneyAtDate(date int64) api.DriverJourney {
	dj := api.NewDriverJourney()
	dj.PassengerPickupDate = date

	return dj
}

func makePassengerJourneyAtDate(date int64) api.PassengerJourney {
	pj := api.NewPassengerJourney()
	pj.PassengerPickupDate = date

	return pj
}

func castDriverToPassenger(p *api.GetDriverJourneysParams) *api.GetPassengerJourneysParams {
	if p == nil {
		return nil
	}

	castedP := api.GetPassengerJourneysParams(*p)
	return &castedP
}

func makeParamsWithDepartureRadius(departureCoord util.Coord, departureRadius float32, driverOrPassenger string) api.GetJourneysParams {
	params := api.NewGetDriverJourneysParams(departureCoord, util.CoordIgnore, 0)
	params.DepartureRadius = &departureRadius

	if driverOrPassenger == "passenger" {
		return castDriverToPassenger(params)
	}

	return params
}

func makeParamsWithArrivalRadius(arrivalCoord util.Coord, arrivalRadius float32, driverOrPassenger string) api.GetJourneysParams {
	params := api.NewGetDriverJourneysParams(util.CoordIgnore, arrivalCoord, 0)
	params.ArrivalRadius = &arrivalRadius

	if driverOrPassenger == "passenger" {
		return castDriverToPassenger(params)
	}

	return params
}

func makeParamsWithTimeDelta(date int, driverOrPassenger string) api.GetJourneysParams {
	params := &api.GetDriverJourneysParams{}
	params.TimeDelta = &date

	if driverOrPassenger == "passenger" {
		return castDriverToPassenger(params)
	}

	return params
}

func makeParamsWithCount(count int, driverOrPassenger string) api.GetJourneysParams {
	params := &api.GetDriverJourneysParams{}
	params.Count = &count

	if driverOrPassenger == "passenger" {
		return castDriverToPassenger(params)
	}

	return params
}

// makeBooking sets status to "WAITING_CONFIRMATION" by default
func makeBooking(bookingID uuid.UUID) *api.Booking {
	return makeBookingWithStatus(bookingID, api.BookingStatusWAITINGCONFIRMATION)
}

func makeBookingWithStatus(bookingID uuid.UUID, status api.BookingStatus) *api.Booking {
	return &api.Booking{Id: bookingID, Status: status}
}

func makeCarpoolBookingEvent(eventID, bookingID uuid.UUID) *api.CarpoolBookingEvent {
	return makeCarpoolBookingEventWithStatus(eventID, bookingID,
		api.BookingStatusWAITINGCONFIRMATION)
}

func makeCarpoolBookingEventWithStatus(eventID, bookingID uuid.UUID, status api.BookingStatus) *api.CarpoolBookingEvent {
	booking := makeBookingWithStatus(bookingID, status)

	carpoolBookingEventData := api.CarpoolBookingEvent_Data{}
	err := carpoolBookingEventData.FromDriverCarpoolBooking(*booking.ToDriverCarpoolBooking())
	panicIf(err)

	carpoolBookingEvent := &api.CarpoolBookingEvent{
		Data:    carpoolBookingEventData,
		Id:      eventID,
		IdToken: "",
	}

	return carpoolBookingEvent
}

// repUUID creates a reproducible UUID
func repUUID(seed int64) uuid.UUID {
	// generate random bytes
	rand.Seed(seed)
	randBytes := make([]byte, 16)
	rand.Read(randBytes)

	uuid, err := uuid.FromBytes(randBytes)
	panicIf(err)

	return uuid
}

// NewBookingsByID populates a BookingsByID with given bookings. It does not
// test if booking is already set.
func NewBookingsByID(bookings ...*api.Booking) BookingsByID {
	var bookingsByID = BookingsByID{}

	for _, booking := range bookings {
		if booking == nil {
			panic(errors.New("attempt to insert nil booking"))
		}

		bookingsByID[booking.Id] = booking
	}

	return bookingsByID
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}
