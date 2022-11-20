package service

import (
	"testing"

	"github.com/fabmob/playground-standard-covoiturage/cmd/api"
	"github.com/fabmob/playground-standard-covoiturage/cmd/test"
)

func TestMustReadDefaultData(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("MustReadDefaultData panics")
		}
	}()
	MustReadDefaultData()
}

func TestDefaultDriverJourneysValidity(t *testing.T) {
	params := requestAll(t, "driver")
	mockDB := NewMockDBWithDefaultData()
	testFun := test.TestGetDriverJourneysResponse

	flags := test.NewFlags()
	flags.DisallowEmpty = true

	testGetJourneysHelper(t, params, mockDB, testFun, flags)
}

func TestDefaultPassengerJourneysValidity(t *testing.T) {
	params := requestAll(t, "passenger")
	mockDB := NewMockDBWithDefaultData()
	testFun := test.TestGetPassengerJourneysResponse

	flags := test.NewFlags()
	flags.DisallowEmpty = true

	testGetJourneysHelper(t, params, mockDB, testFun, flags)
}

func requestAll(t *testing.T, driverOrPassenger string) api.GetJourneysParams {
	t.Helper()

	var (
		largeTimeDelta = int(1e10)
		largeRadius    = float32(1e6)
	)

	switch driverOrPassenger {
	case "driver":
		params := api.GetDriverJourneysParams{}
		params.TimeDelta = &largeTimeDelta
		params.DepartureRadius = &largeRadius
		params.ArrivalRadius = &largeRadius

		return &params

	case "passenger":
		params := api.GetPassengerJourneysParams{}
		params.TimeDelta = &largeTimeDelta
		params.DepartureRadius = &largeRadius
		params.ArrivalRadius = &largeRadius

		return &params

	default:
		panic("invalid driverOrPassenger parameter")
	}
}
