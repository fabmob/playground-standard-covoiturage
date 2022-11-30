package db

import (
	"bytes"
	"encoding/json"
	"io"

	// for the go:embed directive
	_ "embed"

	"github.com/fabmob/playground-standard-covoiturage/cmd/api"
)

type MockDBDataInterface struct {
	DriverJourneys        []api.DriverJourney        `json:"driverJourneys"`
	PassengerJourneys     []api.PassengerJourney     `json:"passengerJourneys"`
	DriverRegularTrips    []api.DriverRegularTrip    `json:"driverRegularTrips"`
	PassengerRegularTrips []api.PassengerRegularTrip `json:"passengerRegularTrips"`
	Bookings              []*api.Booking             `json:"bookings"`
	Users                 []api.User                 `json:"users"`
	Messages              []api.PostMessagesJSONBody `json:"messages"`
}

func toOutputData(m *Mock) MockDBDataInterface {
	outputData := MockDBDataInterface{}

	outputData.DriverJourneys = m.DriverJourneys
	outputData.PassengerJourneys = m.PassengerJourneys
	outputData.DriverRegularTrips = m.DriverRegularTrips
	outputData.PassengerRegularTrips = m.PassengerRegularTrips
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

	_, err = w.Write(jsonData)
	if err != nil {
		return err
	}

	return nil
}

func fromInputData(inputData MockDBDataInterface) *Mock {
	var m = NewMockDB()

	m.DriverJourneys = inputData.DriverJourneys
	m.PassengerJourneys = inputData.PassengerJourneys
	m.DriverRegularTrips = inputData.DriverRegularTrips
	m.PassengerRegularTrips = inputData.PassengerRegularTrips
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
	var data MockDBDataInterface

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
