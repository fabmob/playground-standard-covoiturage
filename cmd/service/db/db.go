package db

import (
	"bytes"
	"fmt"

	// for the go:embed directive
	_ "embed"
	"encoding/json"
	"io"

	"github.com/fabmob/playground-standard-covoiturage/cmd/api"
	"github.com/google/uuid"
)

type DB interface {
	// Getters should never return nil.
	GetDriverJourneys() []api.DriverJourney
	GetPassengerJourneys() []api.PassengerJourney
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
	DriverJourneys    []api.DriverJourney
	PassengerJourneys []api.PassengerJourney
	Bookings          BookingsByID
	Users             []api.User
}

type BookingsByID map[api.BookingId]*api.Booking

// NewMockDB initiates a MockDB with no data
func NewMockDB() *Mock {
	m := Mock{}
	m.DriverJourneys = []api.DriverJourney{}
	m.PassengerJourneys = []api.PassengerJourney{}
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

//////////////////////////////////////////////////////////
// MockDB from data
//////////////////////////////////////////////////////////

type mockDBDataInterface struct {
	DriverJourneys    []api.DriverJourney        `json:"driverJourneys"`
	PassengerJourneys []api.PassengerJourney     `json:"passengerJourneys"`
	Bookings          []*api.Booking             `json:"bookings"`
	Users             []api.User                 `json:"users"`
	Messages          []api.PostMessagesJSONBody `json:"messages"`
}

func toOutputData(m *Mock) mockDBDataInterface {
	outputData := mockDBDataInterface{}

	outputData.DriverJourneys = m.DriverJourneys
	outputData.PassengerJourneys = m.PassengerJourneys
	outputData.Users = m.Users

	outputData.Bookings = make([]*api.Booking, 0, len(m.Bookings))
	for _, booking := range m.Bookings {
		outputData.Bookings = append(outputData.Bookings, booking)
	}

	return outputData
}

func WriteData(m *Mock, w io.Writer) error {
	outputData := toOutputData(m)

	jsonData, err := json.MarshalIndent(outputData, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(jsonData))

	_, err = w.Write(jsonData)
	if err != nil {
		return err
	}

	return nil
}

func fromInputData(inputData mockDBDataInterface) *Mock {
	var m = NewMockDB()

	m.DriverJourneys = inputData.DriverJourneys
	m.PassengerJourneys = inputData.PassengerJourneys
	m.Users = inputData.Users

	m.Bookings = make(BookingsByID, len(inputData.Bookings))

	for _, booking := range inputData.Bookings {
		m.Bookings[booking.Id] = booking
	}

	return m
}

// NewMockDBWithDefaultData initiates a MockDB with default data
func NewMockDBWithDefaultData() *Mock {
	return MustReadDefaultData()
}

// NewMockDBWithData reads journey data from io.Reader with json data.
// It does not validate data against the standard.
func NewMockDBWithData(r io.Reader) (*Mock, error) {
	var data mockDBDataInterface

	bytes, readErr := io.ReadAll(r)
	if readErr != nil {
		return nil, readErr
	}

	err := json.Unmarshal(bytes, &data)

	return fromInputData(data), err
}

// DefaultData stores default json data
//
//go:embed data/defaultData.json
var DefaultData []byte

// MustReadDefaultData reads default data, and panics if any error occurs
func MustReadDefaultData() *Mock {
	mockDB, err := NewMockDBWithData(bytes.NewReader(DefaultData))
	if err != nil {
		panic(err)
	}

	return mockDB
}
