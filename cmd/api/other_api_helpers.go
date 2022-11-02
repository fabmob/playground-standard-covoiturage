package api

import (
	"errors"
	"fmt"
	"net/http"

	"gitlab.com/multi/stdcov-api-test/cmd/util"
)

// ParseGetDriverJourneysOKResponse extracts and parses the data when the
// status is expected to be 200 Status OK. An error is returned if code is not 200.
func ParseGetDriverJourneysOKResponse(response *http.Response) ([]DriverJourney, error) {
	responseObj, err := ParseGetDriverJourneysResponse(response)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected 200 status code, got %d", response.StatusCode)
	}
	if responseObj.JSON200 == nil {
		return nil, errors.New("response with missing data")
	}
	return *responseObj.JSON200, nil

}

// NewGetDriverJourneysParams returns query parameters, looking for a trip
// from "departure" to "arrival" at "departureDate".
func NewGetDriverJourneysParams(
	departure, arrival util.Coord,
	departureDate int,
) *GetDriverJourneysParams {
	return &GetDriverJourneysParams{
		DepartureLat:  float32(departure.Lat),
		DepartureLng:  float32(departure.Lon),
		ArrivalLat:    float32(arrival.Lat),
		ArrivalLng:    float32(arrival.Lon),
		DepartureDate: departureDate,
	}
}

// GetTimeDelta returns the value of TimeDelta if any, or its default value
// otherwise.
func (p GetDriverJourneysParams) GetTimeDelta() int {
	defaultTimeDelta := 900
	if p.TimeDelta == nil {
		return defaultTimeDelta
	}
	return *p.TimeDelta
}

// GetDepartureRadius returns the value of DepartureRadius if any, or its default value
// otherwise.
func (p GetDriverJourneysParams) GetDepartureRadius() float64 {
	defaultDepartureRadius := 1.
	if p.DepartureRadius == nil {
		return defaultDepartureRadius
	}
	return float64(*p.DepartureRadius)
}

// GetArrivalRadius returns the value of ArrivalRadius if any, or its default value
// otherwise.
func (p GetDriverJourneysParams) GetArrivalRadius() float64 {
	defaultArrivalRadius := 1.
	if p.ArrivalRadius == nil {
		return defaultArrivalRadius
	}
	return float64(*p.ArrivalRadius)
}

// NewDriverJourney returns a valid DriverJourney
func NewDriverJourney() DriverJourney {
	dj := DriverJourney{}
	dj.Type = "DYNAMIC"
	dj.Operator = "example.com"
	return dj
}

// NewPassengerJourney returns a valid PassengerJourney
func NewPassengerJourney() PassengerJourney {
	dj := PassengerJourney{}
	departureDate := int64(0)
	dj.DriverDepartureDate = &departureDate
	dj.Type = "DYNAMIC"
	return dj
}
