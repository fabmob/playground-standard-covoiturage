package service

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

// MockDB stores the data of the server in memory
type MockDB struct {
	DriverJourneys    []api.DriverJourney
	PassengerJourneys []api.PassengerJourney
	Bookings          BookingsByID
	Users             []api.User
	Messages          []api.PostMessagesJSONBody
}

type BookingsByID map[uuid.UUID]*api.Booking

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

func (m *MockDB) GetUsers() []api.User {
	if m.Users == nil {
		m.Users = []api.User{}
	}

	return m.Users
}

func (m *MockDB) GetBooking(bookingID uuid.UUID) (*api.Booking, error) {
	bookings := m.GetBookings()

	booking, ok := bookings[bookingID]
	if !ok {
		return nil, MissingBookingErr{}
	}

	return booking, nil
}

// AddBooking adds a new booking to the data. Returns an error if a booking
// with same ID already exists
func (m *MockDB) AddBooking(booking api.Booking) error {
	bookings := m.GetBookings()

	if _, bookingExists := bookings[booking.Id]; bookingExists {
		return fmt.Errorf("booking already exists (ID: %s)", booking.Id)
	}

	bookings[booking.Id] = &booking

	return nil
}

// UpdateBookingStatus updates the status of a booking. Status can only be
// updated for a higher ranked status. If this is not the case, or if the
// booking is not found, returns an error
func (m *MockDB) UpdateBookingStatus(bookingID uuid.UUID, newStatus api.BookingStatus) error {
	booking, err := m.GetBooking(bookingID)
	if err != nil {
		return err
	}

	statusAfter, err := statusIsAfter(newStatus, booking.Status)
	if err != nil {
		return err
	}

	if !statusAfter {
		return StatusAlreadySetErr{}
	}

	booking.Status = newStatus

	return nil
}

type MissingBookingErr struct{}

func (err MissingBookingErr) Error() string {
	return "missing_booking"
}

type StatusAlreadySetErr struct{}

func (err StatusAlreadySetErr) Error() string {
	return "status_already_set"
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

func toOutputData(m *MockDB) mockDBDataInterface {
	outputData := mockDBDataInterface{}

	outputData.DriverJourneys = m.DriverJourneys
	outputData.PassengerJourneys = m.PassengerJourneys
	outputData.Users = m.Users
	outputData.Messages = m.Messages

	outputData.Bookings = make([]*api.Booking, 0, len(m.Bookings))
	for _, booking := range m.Bookings {
		outputData.Bookings = append(outputData.Bookings, booking)
	}

	return outputData
}

func fromInputData(inputData mockDBDataInterface) *MockDB {
	var m = NewMockDB()

	m.DriverJourneys = inputData.DriverJourneys
	m.PassengerJourneys = inputData.PassengerJourneys
	m.Users = inputData.Users
	m.Messages = inputData.Messages

	m.Bookings = make(BookingsByID, len(inputData.Bookings))

	for _, booking := range inputData.Bookings {
		m.Bookings[booking.Id] = booking
	}

	return m
}

// NewMockDBWithDefaultData initiates a MockDB with default data
func NewMockDBWithDefaultData() *MockDB {
	return MustReadDefaultData()
}

// NewMockDBWithData reads journey data from io.Reader with json data.
// It does not validate data against the standard.
func NewMockDBWithData(r io.Reader) (*MockDB, error) {
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
func MustReadDefaultData() *MockDB {
	mockDB, err := NewMockDBWithData(bytes.NewReader(DefaultData))
	if err != nil {
		panic(err)
	}

	return mockDB
}
