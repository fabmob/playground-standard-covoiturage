package db

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMustReadDefaultData(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("MustReadDefaultData panics")
		}
	}()
	MustReadDefaultData()
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

func TestWriteData(t *testing.T) {
	mockDB := NewMockDB()

	var b bytes.Buffer

	err := WriteData(mockDB, &b)
	if err != nil {
		panic(err)
	}

	writtenData, err := io.ReadAll(&b)
	if err != nil {
		panic(err)
	}
	assert.NotEmpty(t, writtenData)
	assert.True(t, isValidJSON(writtenData))
}

func isValidJSON(input []byte) bool {
	var js json.RawMessage
	return json.Unmarshal(input, &js) == nil
}

/* func TestFromInputData(t *testing.T) { */
/* 	inputData := mockDBDataInterface{ */
/* 		DriverJourneys:    makeNDriverJourneys(3), */
/* 		PassengerJourneys: makeNPassengerJourneys(4), */
/* 		Users:             []api.User{makeUser("1", "alice"), makeUser("2", "bob")}, */
/* 		Bookings: []*api.Booking{ */
/* 			makeBooking(repUUID(0)), */
/* 			makeBooking(repUUID(1)), */
/* 		}, */
/* 	} */

/* 	mockDB := fromInputData(inputData) */

/* 	assert.Equal(t, inputData.DriverJourneys, mockDB.DriverJourneys) */
/* 	assert.Equal(t, inputData.PassengerJourneys, mockDB.PassengerJourneys) */
/* 	assert.Equal(t, inputData.Users, mockDB.Users) */
/* 	assert.Equal(t, inputData.Bookings[0], mockDB.Bookings[repUUID(0)]) */
/* 	assert.Equal(t, inputData.Bookings[1], mockDB.Bookings[repUUID(1)]) */
/* } */
