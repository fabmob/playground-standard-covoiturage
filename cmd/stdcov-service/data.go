package main

import (
	"io"

	"gitlab.com/multi/stdcov-api-test/cmd/stdcov-service/server"
)

// ReadJourneyData reads starting journey data from json file
func ReadJourneyData(r io.Reader) ([]server.DriverJourney, error) {

	return nil, nil
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
