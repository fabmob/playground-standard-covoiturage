package service

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"io"

	"gitlab.com/multi/stdcov-api-test/cmd/api"
)

type mockDB struct {
	driverJourneys []api.DriverJourney
}

func NewMockDB() mockDB {
	m := mockDB{}
	m.driverJourneys = []api.DriverJourney{}
	return m
}

func (db *mockDB) PopulateDBWithDefault() error {
	var err error
	db.driverJourneys, err = ReadJourneyData(bytes.NewReader(DriverJourneyJSON))
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
