package main

import (
	"encoding/json"
	"io"
	"os"

	"gitlab.com/multi/stdcov-api-test/cmd/stdcov-service/server"
)

// DriverJourneysData is the in-memory equivalent of the driver journeys
// stored in a database
var DriverJourneysData, _ = ReadJourneyDataFromFile("./data/defaultJourneyData.json")

// ReadJourneyDataFromFile reads a []DriverJourney array from a json file at given
// path
func ReadJourneyDataFromFile(path string) ([]server.DriverJourney, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return ReadJourneyData(f)
}

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
