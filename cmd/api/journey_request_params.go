package api

import "net/http"

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
