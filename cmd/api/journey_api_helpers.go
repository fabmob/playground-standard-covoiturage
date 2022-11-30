package api

import (
	"fmt"

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

// NewTrip returns a valid Trip
func NewTrip() Trip {
	t := Trip{}
	t.Operator = "example.com"

	return t
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

// NewDriverRegularTrip returns a valid DriverRegularTrip
func NewDriverRegularTrip() DriverRegularTrip {
	drt := DriverRegularTrip{}

	return drt
}

// NewPassengerRegularTrip returns a valid PassengerRegularTrip
func NewPassengerRegularTrip() PassengerRegularTrip {
	prt := PassengerRegularTrip{}

	return prt
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

type JourneyOrTripPartialParams interface {
	GetDepartureLat() float64
	GetDepartureLng() float64
	GetDepartureRadius() float64
	GetArrivalLat() float64
	GetArrivalLng() float64
	GetArrivalRadius() float64
	GetTimeDelta() int
	GetCount() *int
}

type GetJourneysParams interface {
	JourneyOrTripPartialParams
	GetDepartureDate() int
}

type GetRegularTripParams interface {
	JourneyOrTripPartialParams
	GetDepartureTimeOfDay() string
	GetDepartureWeekDays() []string
}

var (
	defaultTimeDelta         = 900
	defaultDepartureRadius   = float32(1.)
	defaultArrivalRadius     = float32(1.)
	defaultDepartureWeekdays = []string{"MON", "TUE", "WED", "THU", "FRI", "SAT", "SUN"}
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

///////////////////////////////////////////////
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

///////////////////////////////////////////////
// *GetDriverRegularTripsParams implements GetRegularTripsParams

func (p *GetDriverRegularTripsParams) GetDepartureLat() float64 {
	return float64(p.DepartureLat)
}

func (p *GetDriverRegularTripsParams) GetDepartureLng() float64 {
	return float64(p.DepartureLng)
}

// GetDepartureRadius returns the value of DepartureRadius if any, or its default value
// otherwise.
func (p *GetDriverRegularTripsParams) GetDepartureRadius() float64 {
	return float64(withDefaultValue(p.DepartureRadius, defaultDepartureRadius))
}

func (p *GetDriverRegularTripsParams) GetDepartureTimeOfDay() string {
	return p.DepartureTimeOfDay
}

func (p *GetDriverRegularTripsParams) GetDepartureWeekDays() []string {
	return withDefaultValue(p.DepartureWeekdays, defaultDepartureWeekdays)
}

func (p *GetDriverRegularTripsParams) GetArrivalLat() float64 {
	return float64(p.ArrivalLat)
}

func (p *GetDriverRegularTripsParams) GetArrivalLng() float64 {
	return float64(p.ArrivalLng)
}

// GetArrivalRadius returns the value of ArrivalRadius if any, or its default value
// otherwise.
func (p *GetDriverRegularTripsParams) GetArrivalRadius() float64 {
	return float64(withDefaultValue(p.ArrivalRadius, defaultArrivalRadius))
}

// GetTimeDelta returns the value of TimeDelta if any, or its default value
// otherwise. Implements GetJourneyParams.GetTimeDelta().
func (p *GetDriverRegularTripsParams) GetTimeDelta() int {
	return withDefaultValue(p.TimeDelta, defaultTimeDelta)
}

func (p *GetDriverRegularTripsParams) GetCount() *int {
	return p.Count
}

///////////////////////////////////////////////
// *GetPassengerRegularTripsParams implements GetRegularTripsParams

func (p *GetPassengerRegularTripsParams) GetDepartureLat() float64 {
	return float64(p.DepartureLat)
}

func (p *GetPassengerRegularTripsParams) GetDepartureLng() float64 {
	return float64(p.DepartureLng)
}

// GetDepartureRadius returns the value of DepartureRadius if any, or its default value
// otherwise.
func (p *GetPassengerRegularTripsParams) GetDepartureRadius() float64 {
	return float64(withDefaultValue(p.DepartureRadius, defaultDepartureRadius))
}

func (p *GetPassengerRegularTripsParams) GetDepartureTimeOfDay() string {
	return p.DepartureTimeOfDay
}

func (p *GetPassengerRegularTripsParams) GetDepartureWeekDays() []string {
	return withDefaultValue(p.DepartureWeekdays, defaultDepartureWeekdays)
}

func (p *GetPassengerRegularTripsParams) GetArrivalLat() float64 {
	return float64(p.ArrivalLat)
}

func (p *GetPassengerRegularTripsParams) GetArrivalLng() float64 {
	return float64(p.ArrivalLng)
}

// GetArrivalRadius returns the value of ArrivalRadius if any, or its default value
// otherwise.
func (p *GetPassengerRegularTripsParams) GetArrivalRadius() float64 {
	return float64(withDefaultValue(p.ArrivalRadius, defaultArrivalRadius))
}

// GetTimeDelta returns the value of TimeDelta if any, or its default value
// otherwise. Implements GetJourneyParams.GetTimeDelta().
func (p *GetPassengerRegularTripsParams) GetTimeDelta() int {
	return withDefaultValue(p.TimeDelta, defaultTimeDelta)
}

func (p *GetPassengerRegularTripsParams) GetCount() *int {
	return p.Count
}

////////////////////////////////////////////////////
// Helpers
////////////////////////////////////////////////////

// withDefaultValue takes a pointer, and returns the value pointed at, or a
// default value if the pointer is nil
func withDefaultValue[T any](t *T, defaultValue T) T {
	if t == nil {
		return defaultValue
	}

	return *t
}
