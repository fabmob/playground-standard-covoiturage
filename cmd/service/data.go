package service

import (
	"bytes"
	// for the go:embed directive
	_ "embed"
	"encoding/json"
	"io"

	"gitlab.com/multi/stdcov-api-test/cmd/api"
)

// MockDB stores the data of the server in memory
type MockDB struct {
	DriverJourneys    []api.DriverJourney    `json:"driverJourneys"`
	PassengerJourneys []api.PassengerJourney `json:"passengerJourneys"`
}

// NewMockDB initiates a MockDB with no data
func NewMockDB() *MockDB {
	m := MockDB{}
	m.DriverJourneys = []api.DriverJourney{}
	m.PassengerJourneys = []api.PassengerJourney{}
	return &m
}

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
