package main

import (
	"bytes"
	"testing"
)

func TestDefaultJourneyData(t *testing.T) {
	_, err := ReadJourneyData(bytes.NewReader(DriverJourneyJSON))
	if err != nil {
		t.Error(err)
	}
}
