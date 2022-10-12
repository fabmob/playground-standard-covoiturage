package main

import "gitlab.com/multi/stdcov-api-test/cmd/stdcov-service/server"

var journeys = []server.DriverJourney{
	{
		Driver: server.User{
			// User's alias.
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
