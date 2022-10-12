package main

import (
	"strings"
	"testing"

	"gitlab.com/multi/stdcov-api-test/cmd/stdcov-service/server"
)

func TestReadJourneyData(t *testing.T) {
	expected := []server.DriverJourney{
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

	jsonDriverJourney := `[
	{
		"driver": {
			"alias": "bob",
			"id":    "1"
		},
		"operator":            "operator.example.org",
		"duration":            3600,
		"passengerDropLat":    48.8450234,
		"passengerDropLng":    2.3997529,
		"passengerPickupDate": 1665579951,
		"passengerPickupLat":  47.461737,
		"type": "DYNAMIC"
	}
]`
	r := strings.NewReader(jsonDriverJourney)
	got, err := ReadJourneyData(r)

	if err != nil {
		t.Error("Error while reading test string")
	}
	if len(got) != len(expected) {
		t.Logf("Expected length of data: %d", len(expected))
		t.Logf("Got length of data: %d", len(got))
		t.Error()
	}
}
