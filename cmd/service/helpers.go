package service

import (
	"fmt"
	"math"
	"time"

	"github.com/fabmob/playground-standard-covoiturage/cmd/api"
	"github.com/fabmob/playground-standard-covoiturage/cmd/service/db"
	"github.com/fabmob/playground-standard-covoiturage/cmd/util"
	"github.com/google/uuid"
)

// keepNFirst keeps n first elements of slice, or returns the slice untouched
// if its length is inferior to n
func keepNFirst[K any](slice []K, n int) []K {
	if len(slice) > n {
		return slice[0:n]
	}

	return slice
}

var statusToIntMap = map[api.BookingStatus]int{
	api.BookingStatusWAITINGCONFIRMATION:        0,
	api.BookingStatusCONFIRMED:                  1,
	api.BookingStatusCOMPLETEDPENDINGVALIDATION: 2,
	api.BookingStatusVALIDATED:                  3,
	api.BookingStatusCANCELLED:                  4,
}

// statusIsAfter checks if status1 is strictly after status2
func statusIsAfter(status1, status2 api.BookingStatus) (bool, error) {
	status1Rank, err := statusRank(status1)
	if err != nil {
		return false, err
	}

	status2Rank, err := statusRank(status2)
	if err != nil {
		return false, err
	}

	return status1Rank > status2Rank, nil
}

type StatusAlreadySetErr struct{}

func (err StatusAlreadySetErr) Error() string {
	return "status_already_set"
}

// UpdateBookingStatus updates the status of a booking. Status can only be
// updated for a higher ranked status. If this is not the case, or if the
// booking is not found, returns an error
func UpdateBookingStatus(m db.DB, bookingID uuid.UUID, newStatus api.BookingStatus) error {
	booking, err := m.GetBooking(bookingID)
	if err != nil {
		return err
	}

	statusAfter, err := statusIsAfter(newStatus, booking.Status)
	if err != nil {
		return err
	}

	if !statusAfter {
		return StatusAlreadySetErr{}
	}

	booking.Status = newStatus

	return nil
}

func statusRank(status api.BookingStatus) (int, error) {
	statusRank, ok2 := statusToIntMap[status]
	if !ok2 {
		return 0, fmt.Errorf("%s is not a valid status", status)
	}

	return statusRank, nil
}

// errorBody creates an api.BadRequest body from a go error
func errorBody(err error) api.BadRequest {
	errStr := err.Error()
	return api.BadRequest{Error: &errStr}
}

func userExists(user api.User, users []api.User) bool {
	for _, existingUser := range users {
		if existingUser.Id == user.Id &&
			existingUser.Operator == user.Operator {
			return true
		}
	}

	return false
}

// keepJourney checks a journey's trip and schedule (common properties to
// passenger and driver journeys)  against the query parameters.
func keepJourney(params api.GetJourneysParams, trip api.Trip, schedule api.JourneySchedule) bool {
	tripOK := keepTrip(params, trip)

	var (
		expectedDate = float64(params.GetDepartureDate())
		gotDate      = float64(schedule.PassengerPickupDate)
		maxDiff      = float64(params.GetTimeDelta())
	)

	timeDeltaOK := math.Abs(gotDate-expectedDate) < maxDiff

	return tripOK && timeDeltaOK
}

// filterSchedules filters schedules for regular trips given query parameters.
// Expects non-nil "schedules" argument.
func keepSchedule(params api.GetRegularTripParams, schedule api.Schedule) (bool, error) {

	if schedule.PassengerPickupDay == nil || schedule.PassengerPickupTimeOfDay == nil {
		return false, nil
	}

	passengerPickupDay := *schedule.PassengerPickupDay
	passengerPickupTimeOfDay := *schedule.PassengerPickupTimeOfDay

	validWeekDay := isAllowedWeekday(passengerPickupDay, params.GetDepartureWeekDays())

	d, err := durationBetweenTimeOfDays(passengerPickupTimeOfDay, params.GetDepartureTimeOfDay())
	if err != nil {
		return false, err
	}
	validTimeOfDay := int(d) <= params.GetTimeDelta()

	validPeriod := true

	if schedule.JourneySchedules != nil {
		validPeriod = anyJourneyScheduleInMinMax(*schedule.JourneySchedules,
			params.GetMinDepartureDate(), params.GetMaxDepartureDate(),
		)
	}

	return validWeekDay && validTimeOfDay && validPeriod, nil
}

// isAllowedWeekday checks if day is in allowedDays
func isAllowedWeekday(day api.SchedulePassengerPickupDay, allowedDays []string) bool {
	validDay := false

	for _, allowedDay := range allowedDays {
		if string(day) == allowedDay {
			validDay = true
			break
		}
	}

	return validDay
}

// anyJourneyScheduleInMinMax checks if any JourneySchedule has a UNIX
// passenger pickup date in between min and max UNIX dates.
func anyJourneyScheduleInMinMax(journeySchedules []api.JourneySchedule, min, max *int) bool {
	for _, js := range journeySchedules {
		if min != nil {
			if belowMin := js.PassengerPickupDate < int64(*min); belowMin {
				continue
			}
		}

		if max != nil {
			if aboveMax := js.PassengerPickupDate > int64(*max); aboveMax {
				continue
			}
		}

		return true
	}

	return false
}

// durationBetweenTimeOfDays returns the duration between two partial times,
// considering them on the same day.
func durationBetweenTimeOfDays(t1, t2 string) (float64, error) {
	time1, err := time.Parse("15:04:05", t1)
	if err != nil {
		return 0, err
	}

	time2, err := time.Parse("15:04:05", t2)
	if err != nil {
		return 0, err
	}

	return time1.Sub(time2).Abs().Seconds(), nil
}

// keepTrip checks if a trip object is compliant with the query parameters
func keepTrip(params api.JourneyOrTripPartialParams, trip api.Trip) bool {
	coordsRequestDeparture := util.Coord{
		Lat: float64(params.GetDepartureLat()),
		Lon: float64(params.GetDepartureLng()),
	}
	coordsResponseDeparture := util.Coord{
		Lat: trip.PassengerPickupLat,
		Lon: trip.PassengerPickupLng,
	}
	departureRadiusOK := util.Distance(coordsRequestDeparture, coordsResponseDeparture) <=
		params.GetDepartureRadius()

	coordsRequestArrival := util.Coord{
		Lat: float64(params.GetArrivalLat()),
		Lon: float64(params.GetArrivalLng()),
	}
	coordsResponseArrival := util.Coord{
		Lat: trip.PassengerDropLat,
		Lon: trip.PassengerDropLng,
	}
	arrivalRadiusOK := util.Distance(coordsRequestArrival, coordsResponseArrival) <=
		params.GetArrivalRadius()

	return arrivalRadiusOK && departureRadiusOK
}
