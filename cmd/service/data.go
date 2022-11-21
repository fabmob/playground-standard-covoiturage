package service

import (
	"bytes"
	"errors"

	// for the go:embed directive
	_ "embed"
	"encoding/json"
	"io"

	"github.com/fabmob/playground-standard-covoiturage/cmd/api"
	"github.com/google/uuid"
)

// MockDB stores the data of the server in memory
type MockDB struct {
	DriverJourneys    []api.DriverJourney    `json:"driverJourneys"`
	PassengerJourneys []api.PassengerJourney `json:"passengerJourneys"`
	Bookings          BookingsByID           `json:"bookings"`
}

// NewMockDB initiates a MockDB with no data
func NewMockDB() *MockDB {
	m := MockDB{}
	m.DriverJourneys = []api.DriverJourney{}
	m.PassengerJourneys = []api.PassengerJourney{}
	m.Bookings = BookingsByID{}

	return &m
}

func (m *MockDB) GetDriverJourneys() []api.DriverJourney {
	if m.DriverJourneys == nil {
		m.DriverJourneys = []api.DriverJourney{}
	}

	return m.DriverJourneys
}

func (m *MockDB) GetPassengerJourneys() []api.PassengerJourney {
	if m.PassengerJourneys == nil {
		m.PassengerJourneys = []api.PassengerJourney{}
	}

	return m.PassengerJourneys
}

func (m *MockDB) GetBookings() BookingsByID {
	if m.Bookings == nil {
		m.Bookings = BookingsByID{}
	}

	return m.Bookings
}

func (m *MockDB) GetBooking(bookingID uuid.UUID) (*api.Booking, error) {
	bookings := m.GetBookings()

	booking, ok := bookings[bookingID]
	if !ok {
		return nil, errors.New("missing_booking")
	}

	return booking, nil
}

// AddBooking adds a new booking to the data. Returns an error if a booking
// with same ID already exists
func (m *MockDB) AddBooking(booking api.Booking) error {
	bookings := m.GetBookings()

	if _, bookingExists := bookings[booking.Id]; bookingExists {
		return errors.New("booking already exists")
	}

	bookings[booking.Id] = &booking

	return nil
}

type BookingsByID map[uuid.UUID]*api.Booking

// NewMockDBWithDefaultData initiates a MockDB with default data
func NewMockDBWithDefaultData() *MockDB {
	return MustReadDefaultData()
}

// JSONData stores default driver journey json data
//
//go:embed data/defaultData.json
var JSONData []byte

// DriverJourneysData is the in-memory equivalent of the driver journeys
// stored in a database

// MustReadDefaultData reads default data, and panics if any error occurs
func MustReadDefaultData() *MockDB {
	mockDB, err := ReadData(bytes.NewReader(JSONData))
	if err != nil {
		panic(err)
	}

	return mockDB
}

// ReadData reads journey data from io.Reader with json data.
// It does not validate data.
func ReadData(r io.Reader) (*MockDB, error) {
	var data MockDB

	bytes, readErr := io.ReadAll(r)
	if readErr != nil {
		return nil, readErr
	}

	err := json.Unmarshal(bytes, &data)

	return &data, err
}
