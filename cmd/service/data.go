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
	DriverJourneys    []api.DriverJourney
	PassengerJourneys []api.PassengerJourney
}

// NewMockDB initiates a MockDB with no data
func NewMockDB() MockDB {
	m := MockDB{}
	m.DriverJourneys = []api.DriverJourney{}
	return m
}

// PopulateDBWithDefault populates the MockDB with default data
func (db *MockDB) PopulateDBWithDefault() error {
	var err error
	db.DriverJourneys, err = ReadJourneyData(bytes.NewReader(DriverJourneyJSON))
	return err
}

// DriverJourneyJSON stores default driver journey json data
//
//go:embed data/defaultJourneyData.json
var DriverJourneyJSON []byte

// DriverJourneysData is the in-memory equivalent of the driver journeys
// stored in a database

// ReadJourneyData reads journey data from io.Reader with json data
// It does not validate data
func ReadJourneyData(r io.Reader) ([]api.DriverJourney, error) {
	var journeyData []api.DriverJourney
	bytes, readErr := io.ReadAll(r)
	if readErr != nil {
		return nil, readErr
	}

	err := json.Unmarshal(bytes, &journeyData)
	return journeyData, err
}
