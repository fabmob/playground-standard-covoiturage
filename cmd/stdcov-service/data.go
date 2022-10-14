package main

import (
	"encoding/json"
	"errors"
	"io"
	"os"

	"gitlab.com/multi/stdcov-api-test/cmd/stdcov-service/server"
)

func ReadJourneyDataFromFile(path string) ([]server.DriverJourney, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return ReadJourneyData(f)
}

// ReadJourneyData reads starting journey data from json file
func ReadJourneyData(r io.Reader) ([]server.DriverJourney, error) {
	var journeyData []server.DriverJourney
	bytes, readErr := io.ReadAll(r)
	if readErr != nil {
		return nil, readErr
	}

	err := json.Unmarshal(bytes, &journeyData)
	if journeyData == nil {
		return nil, errors.New("no journey data to parse")
	}
	return journeyData, err
}

var journeys = []server.DriverJourney{
	{
		Driver: server.User{
			Alias: "bob",
			Id:    "1",
		},
		Operator:            "operator.example.org",
		Duration:            3600,
		PassengerDropLat:    48.8450234,
		PassengerDropLng:    2.3997529,
		PassengerPickupDate: 1665579951,
		PassengerPickupLat:  47.461737,
		Type:                server.DriverJourneyTypeDYNAMIC,
	},
}
