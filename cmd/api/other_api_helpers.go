package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
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

// NewGetPassengerJourneysParams returns query parameters, looking for a trip
// from "departure" to "arrival" at "departureDate".
func NewGetPassengerJourneysParams(
	departure, arrival util.Coord,
	departureDate int,
) *GetPassengerJourneysParams {

	return &GetPassengerJourneysParams{
		DepartureLat:  float32(departure.Lat),
		DepartureLng:  float32(departure.Lon),
		ArrivalLat:    float32(arrival.Lat),
		ArrivalLng:    float32(arrival.Lon),
		DepartureDate: departureDate,
	}
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
	pj := PassengerJourney{}
	departureDate := int64(0)
	pj.Operator = "example.com"
	pj.DriverDepartureDate = &departureDate
	pj.Type = "DYNAMIC"

	return pj
}

func GetJourneys(s ServerInterface, ctx echo.Context, params GetJourneysParams) error {
	switch v := params.(type) {
	case *GetPassengerJourneysParams:
		return s.GetPassengerJourneys(ctx, *params.(*GetPassengerJourneysParams))

	case *GetDriverJourneysParams:
		return s.GetDriverJourneys(ctx, *params.(*GetDriverJourneysParams))

	default:
		return fmt.Errorf("unknown journey type %v: only GetDriverJourneysParams and GetPassengerJourneys supported", v)
	}
}
