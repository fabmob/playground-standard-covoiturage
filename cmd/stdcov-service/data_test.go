package main

import (
	"fmt"
	"testing"
)

func TestDefaultJourneyData(t *testing.T) {
	defaultDataFile := "./data/defaultJourneyData.json"
	v, err := ReadJourneyDataFromFile(defaultDataFile)
	fmt.Printf("%+v\n", v)
	if err != nil {
		t.Error(err)
	}
}
