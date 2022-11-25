package db

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"

	"github.com/fabmob/playground-standard-covoiturage/cmd/api"
	"github.com/google/uuid"
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

func TestNewMockDBWithData(t *testing.T) {

	var (
		b bytes.Buffer

		data = mockDBDataInterface{}

		id      = uuid.New()
		booking = api.Booking{Id: id}
	)

	data.Bookings = []*api.Booking{&booking}

	bookingBytes, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	_, err = b.Write(bookingBytes)
	if err != nil {
		panic(err)
	}

	mockDB, err := NewMockDBWithData(&b)

	assert.Nil(t, err)

	_, ok := mockDB.Bookings[id]
	assert.True(t, ok)
}
