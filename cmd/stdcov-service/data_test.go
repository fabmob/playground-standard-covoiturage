package main

import (
	"testing"
)

func TestDefaultJourneyData(t *testing.T) {
	defaultDataFile := "./data/defaultJourneyData.json"
	_, err := ReadJourneyDataFromFile(defaultDataFile)
	if err != nil {
		t.Error(err)
	}
}
