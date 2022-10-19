package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"io"

	"gitlab.com/multi/stdcov-api-test/cmd/stdcov-service/server"
)

// DriverJourneyJSON stores default driver journey json data
//go:embed data/defaultJourneyData.json
var DriverJourneyJSON []byte

// DriverJourneysData is the in-memory equivalent of the driver journeys
// stored in a database
var DriverJourneysData, _ = ReadJourneyData(bytes.NewReader(DriverJourneyJSON))

// ReadJourneyData reads journey data from io.Reader with json data
// It does not validate data
func ReadJourneyData(r io.Reader) ([]server.DriverJourney, error) {
	var journeyData []server.DriverJourney
	bytes, readErr := io.ReadAll(r)
	if readErr != nil {
		return nil, readErr
	}

	err := json.Unmarshal(bytes, &journeyData)
	return journeyData, err
}
