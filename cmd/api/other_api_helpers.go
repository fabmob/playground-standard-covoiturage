package api

import (
	"errors"
	"fmt"
	"net/http"
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

func NewGetDriverJourneysParams(
	departureLat, departureLng, arrivalLat, arrivalLng float32,
	departureDate int,
) *GetDriverJourneysParams {
	defaultTimeDelta := 900
	defaultDepartureRadius := float32(1.)
	defaultArrivalRadius := float32(1.)
	return &GetDriverJourneysParams{
		DepartureLat:    departureLat,
		DepartureLng:    departureLng,
		ArrivalLat:      arrivalLat,
		ArrivalLng:      arrivalLng,
		DepartureDate:   departureDate,
		TimeDelta:       &defaultTimeDelta,
		DepartureRadius: &defaultDepartureRadius,
		ArrivalRadius:   &defaultArrivalRadius,
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
