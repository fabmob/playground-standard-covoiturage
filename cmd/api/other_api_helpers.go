package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/fabmob/playground-standard-covoiturage/cmd/util"
	"github.com/labstack/echo/v4"
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

// ToBooking performs a lossy conversion from driverCarpoolBooking to Booking
// objects
func (dcb DriverCarpoolBooking) ToBooking() *Booking {
	return &Booking{
		Id:                     dcb.Id,
		Driver:                 dcb.Driver,
		Passenger:              User{},
		PassengerPickupLat:     dcb.PassengerPickupLat,
		PassengerPickupLng:     dcb.PassengerPickupLng,
		PassengerDropLat:       dcb.PassengerDropLat,
		PassengerDropLng:       dcb.PassengerDropLng,
		PassengerPickupAddress: dcb.PassengerPickupAddress,
		PassengerDropAddress:   dcb.PassengerDropAddress,
		Status:                 BookingStatus(dcb.Status),
		Duration:               dcb.Duration,
		Distance:               dcb.Distance,
		WebUrl:                 &dcb.WebUrl,
		Car:                    dcb.Car,
		Price:                  dcb.Price,
	}
}

// ToBooking performs a lossy conversion from passengerCarpoolBooking to
// Booking objects
func (pcb PassengerCarpoolBooking) ToBooking() *Booking {
	return &Booking{
		Id:                     pcb.Id,
		Driver:                 User{},
		Passenger:              pcb.Passenger,
		PassengerPickupLat:     pcb.PassengerPickupLat,
		PassengerPickupLng:     pcb.PassengerPickupLng,
		PassengerDropLat:       pcb.PassengerDropLat,
		PassengerDropLng:       pcb.PassengerDropLng,
		PassengerPickupAddress: pcb.PassengerPickupAddress,
		PassengerDropAddress:   pcb.PassengerDropAddress,
		Status:                 BookingStatus(pcb.Status),
		Duration:               pcb.Duration,
		Distance:               pcb.Distance,
		WebUrl:                 &pcb.WebUrl,
	}
}

// ToDriverCarpoolBooking performs a lossy conversion from
// Booking to DriverCarpoolBooking objects
func (b Booking) ToDriverCarpoolBooking() *DriverCarpoolBooking {
	if b.WebUrl == nil {
		mockURL := ""
		b.WebUrl = &mockURL
	}
	return &DriverCarpoolBooking{
		Car:            b.Car,
		Driver:         b.Driver,
		Price:          b.Price,
		CarpoolBooking: b.toCarpoolBooking(),
	}
}

// ToPassengerCarpoolBooking performs a lossy conversion from
// Booking to PassengerCarpoolBooking objects
func (b Booking) ToPassengerCarpoolBooking() *PassengerCarpoolBooking {
	return &PassengerCarpoolBooking{
		Passenger:      b.Passenger,
		CarpoolBooking: b.toCarpoolBooking(),
	}
}

func (b Booking) toCarpoolBooking() CarpoolBooking {
	if b.WebUrl == nil {
		mockURL := ""
		b.WebUrl = &mockURL
	}
	return CarpoolBooking{
		Distance:               b.Distance,
		Duration:               b.Duration,
		Id:                     b.Id,
		PassengerDropAddress:   b.PassengerDropAddress,
		PassengerDropLat:       b.PassengerDropLat,
		PassengerDropLng:       b.PassengerDropLng,
		PassengerPickupAddress: b.PassengerPickupAddress,
		PassengerPickupDate:    b.PassengerPickupDate,
		PassengerPickupLat:     b.PassengerPickupLat,
		PassengerPickupLng:     b.PassengerPickupLng,
		Status:                 CarpoolBookingStatus(b.Status),
		WebUrl:                 *b.WebUrl,
	}
}
