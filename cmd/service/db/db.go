// Package db  handles data storage and manipulation for the API server.
//
// It exports type `Mock` used to store data in memory, but it can be replaced
// with another storage with the interface `DB` (which `Mock` implements).
//
// A MockDB can be initialized with data, given the data is in json format as
// expected by `MockDBDataInterface`. It can also write its data in json
// format (through the same `MockDBDataInterface`) with the function `WriteData`.
package db

import (
	"fmt"

	"github.com/fabmob/playground-standard-covoiturage/cmd/api"
	"github.com/google/uuid"
)

type DB interface {
	// Getters should never return nil.
	GetDriverJourneys() []api.DriverJourney
	GetPassengerJourneys() []api.PassengerJourney
	GetDriverRegularTrips() []api.DriverRegularTrip
	GetPassengerRegularTrips() []api.PassengerRegularTrip

	GetUsers() []api.User
	GetBookings() BookingsByID

	// GetBooking should return a MissingBookingErr if not found
	GetBooking(api.BookingId) (*api.Booking, error)

	// AddBooking adds a booking to the db, but fails if a booking with same ID
	// already exists
	AddBooking(api.Booking) error
}

// Mock stores the data of the server in memory
type Mock struct {
	DriverJourneys        []api.DriverJourney
	PassengerJourneys     []api.PassengerJourney
	DriverRegularTrips    []api.DriverRegularTrip
	PassengerRegularTrips []api.PassengerRegularTrip
	Bookings              BookingsByID
	Users                 []api.User
}

type BookingsByID map[api.BookingId]*api.Booking

// NewMockDB initiates a MockDB with no data
func NewMockDB() *Mock {
	m := Mock{}
	m.DriverJourneys = []api.DriverJourney{}
	m.PassengerJourneys = []api.PassengerJourney{}
	m.DriverRegularTrips = []api.DriverRegularTrip{}
	m.PassengerRegularTrips = []api.PassengerRegularTrip{}
	m.Bookings = BookingsByID{}
	m.Users = []api.User{}

	return &m
}

func (m *Mock) GetDriverJourneys() []api.DriverJourney {
	if m.DriverJourneys == nil {
		m.DriverJourneys = []api.DriverJourney{}
	}

	return m.DriverJourneys
}

func (m *Mock) GetPassengerJourneys() []api.PassengerJourney {
	if m.PassengerJourneys == nil {
		m.PassengerJourneys = []api.PassengerJourney{}
	}

	return m.PassengerJourneys
}

func (m *Mock) GetDriverRegularTrips() []api.DriverRegularTrip {
	if m.DriverRegularTrips == nil {
		m.DriverRegularTrips = []api.DriverRegularTrip{}
	}

	return m.DriverRegularTrips
}

func (m *Mock) GetPassengerRegularTrips() []api.PassengerRegularTrip {
	if m.PassengerRegularTrips == nil {
		m.PassengerRegularTrips = []api.PassengerRegularTrip{}
	}

	return m.PassengerRegularTrips
}

func (m *Mock) GetBookings() BookingsByID {
	if m.Bookings == nil {
		m.Bookings = BookingsByID{}
	}

	return m.Bookings
}

func (m *Mock) GetUsers() []api.User {
	if m.Users == nil {
		m.Users = []api.User{}
	}

	return m.Users
}

func (m *Mock) GetBooking(bookingID uuid.UUID) (*api.Booking, error) {
	bookings := m.GetBookings()

	booking, ok := bookings[bookingID]
	if !ok {
		return nil, MissingBookingErr{}
	}

	return booking, nil
}

// AddBooking adds a new booking to the data. Returns an error if a booking
// with same ID already exists
func (m *Mock) AddBooking(booking api.Booking) error {
	bookings := m.GetBookings()

	if _, bookingExists := bookings[booking.Id]; bookingExists {
		return fmt.Errorf("booking already exists (ID: %s)", booking.Id)
	}

	bookings[booking.Id] = &booking

	return nil
}

type MissingBookingErr struct{}

func (err MissingBookingErr) Error() string {
	return "missing_booking"
}
