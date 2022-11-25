package api

import (
	"fmt"
	"net/http"

	"github.com/fabmob/playground-standard-covoiturage/cmd/util"
	"github.com/labstack/echo/v4"
)

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

///////////////////////////////////////////////

type RequestParams interface {
	MakeRequest(server string) (*http.Request, error)
}

func (p *GetDriverJourneysParams) MakeRequest(server string) (*http.Request, error) {
	return NewGetDriverJourneysRequest(server, p)
}

func (p *GetPassengerJourneysParams) MakeRequest(server string) (*http.Request, error) {
	return NewGetPassengerJourneysRequest(server, p)
}

type GetJourneysParams interface {
	RequestParams
	GetDepartureLat() float64
	GetDepartureLng() float64
	GetDepartureRadius() float64
	GetDepartureDate() int
	GetArrivalLat() float64
	GetArrivalLng() float64
	GetArrivalRadius() float64
	GetTimeDelta() int
	GetCount() *int
}

var (
	defaultTimeDelta       = 900
	defaultDepartureRadius = float32(1.)
	defaultArrivalRadius   = float32(1.)
)

// *GetDriverJourneysParams implements GetJourneysParams

func (p *GetDriverJourneysParams) GetDepartureLat() float64 {
	return float64(p.DepartureLat)
}

func (p *GetDriverJourneysParams) GetDepartureLng() float64 {
	return float64(p.DepartureLng)
}

// GetDepartureRadius returns the value of DepartureRadius if any, or its default value
// otherwise.
func (p *GetDriverJourneysParams) GetDepartureRadius() float64 {
	return float64(withDefaultValue(p.DepartureRadius, defaultDepartureRadius))
}

func (p *GetDriverJourneysParams) GetDepartureDate() int {
	return p.DepartureDate
}

func (p *GetDriverJourneysParams) GetArrivalLat() float64 {
	return float64(p.ArrivalLat)
}

func (p *GetDriverJourneysParams) GetArrivalLng() float64 {
	return float64(p.ArrivalLng)
}

// GetArrivalRadius returns the value of ArrivalRadius if any, or its default value
// otherwise.
func (p *GetDriverJourneysParams) GetArrivalRadius() float64 {
	return float64(withDefaultValue(p.ArrivalRadius, defaultArrivalRadius))
}

// GetTimeDelta returns the value of TimeDelta if any, or its default value
// otherwise. Implements GetJourneyParams.GetTimeDelta().
func (p *GetDriverJourneysParams) GetTimeDelta() int {
	return withDefaultValue(p.TimeDelta, defaultTimeDelta)
}

func (p *GetDriverJourneysParams) GetCount() *int {
	return p.Count
}

// *GetPassengerJourneysParams implements GetJourneysParams

func (p *GetPassengerJourneysParams) GetDepartureLat() float64 {
	return float64(p.DepartureLat)
}

func (p *GetPassengerJourneysParams) GetDepartureLng() float64 {
	return float64(p.DepartureLng)
}

// GetDepartureRadius returns the value of DepartureRadius if any, or its default value
// otherwise.
func (p *GetPassengerJourneysParams) GetDepartureRadius() float64 {
	return float64(withDefaultValue(p.DepartureRadius, defaultDepartureRadius))
}

func (p *GetPassengerJourneysParams) GetDepartureDate() int {
	return p.DepartureDate
}

func (p *GetPassengerJourneysParams) GetArrivalLat() float64 {
	return float64(p.ArrivalLat)
}

func (p *GetPassengerJourneysParams) GetArrivalLng() float64 {
	return float64(p.ArrivalLng)
}

// GetArrivalRadius returns the value of ArrivalRadius if any, or its default value
// otherwise.
func (p *GetPassengerJourneysParams) GetArrivalRadius() float64 {
	return float64(withDefaultValue(p.ArrivalRadius, defaultArrivalRadius))
}

// GetTimeDelta returns the value of TimeDelta if any, or its default value
// otherwise. Implements GetJourneyParams.GetTimeDelta().
func (p *GetPassengerJourneysParams) GetTimeDelta() int {
	return withDefaultValue(p.TimeDelta, defaultTimeDelta)
}

func (p *GetPassengerJourneysParams) GetCount() *int {
	return p.Count
}

////////////////////////////////////////////////////
// Helpers
////////////////////////////////////////////////////

// withDefaultValue takes a pointer, and returns the value pointed at, or a
// default value if the pointer is nil
func withDefaultValue[T int | float32 | float64](t *T, defaultValue T) T {
	if t == nil {
		return defaultValue
	}

	return *t
}
