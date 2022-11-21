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

func TestMockDB_GetBookings(t *testing.T) {

	t.Run("GetBookings is non-nil even if bookings is nil", func(t *testing.T) {
		db := NewMockDB()
		db.Bookings = nil

		if db.GetBookings() == nil {
			t.Error("GetBookings should never return nil")
		}
	})

	t.Run("GetBookings initialize `Bookings` property as a side effect", func(t *testing.T) {
		db := NewMockDB()
		db.Bookings = nil

		_ = db.GetBookings()

		if db.Bookings == nil {
			t.Error("GetBookings should have as side effect to initialize `Bookings` property")
		}
	})
}
