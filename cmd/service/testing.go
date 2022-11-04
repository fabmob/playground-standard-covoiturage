package service

import (
	"github.com/fabmob/playground-standard-covoiturage/cmd/api"
	"github.com/fabmob/playground-standard-covoiturage/cmd/util"
)

func makeNDriverJourneys(n int) []api.DriverJourney {
	driverJourneys := make([]api.DriverJourney, 0, n)

	for i := 0; i < n; i++ {
		driverJourneys = append(driverJourneys, api.NewDriverJourney())
	}

	return driverJourneys
}

func makeNPassengerJourneys(n int) []api.PassengerJourney {
	passengerJourneys := make([]api.PassengerJourney, 0, n)

	for i := 0; i < n; i++ {
		passengerJourneys = append(passengerJourneys, api.NewPassengerJourney())
	}

	return passengerJourneys
}

func makeDriverJourneyAtCoords(coordPickup, coordDrop util.Coord) api.DriverJourney {
	dj := api.NewDriverJourney()
	updateTripCoords(&dj.Trip, coordPickup, coordDrop)

	return dj
}

func makePassengerJourneyAtCoords(coordPickup, coordDrop util.Coord) api.PassengerJourney {
	pj := api.NewPassengerJourney()
	updateTripCoords(&pj.Trip, coordPickup, coordDrop)

	return pj
}

func updateTripCoords(t *api.Trip, coordPickup, coordDrop util.Coord) {
	t.PassengerPickupLat = coordPickup.Lat
	t.PassengerPickupLng = coordPickup.Lon
	t.PassengerDropLat = coordDrop.Lat
	t.PassengerDropLng = coordDrop.Lon
}

func makeDriverJourneyAtDate(date int64) api.DriverJourney {
	dj := api.NewDriverJourney()
	dj.PassengerPickupDate = date

	return dj
}

func makePassengerJourneyAtDate(date int64) api.PassengerJourney {
	pj := api.NewPassengerJourney()
	pj.PassengerPickupDate = date

	return pj
}

func castDriverToPassenger(p *api.GetDriverJourneysParams) *api.GetPassengerJourneysParams {
	if p == nil {
		return nil
	}

	castedP := api.GetPassengerJourneysParams(*p)
	return &castedP
}

func makeParamsWithDepartureRadius(departureCoord util.Coord, departureRadius float32, driverOrPassenger string) api.GetJourneysParams {
	params := api.NewGetDriverJourneysParams(departureCoord, util.CoordIgnore, 0)
	params.DepartureRadius = &departureRadius

	if driverOrPassenger == "passenger" {
		return castDriverToPassenger(params)
	}

	return params
}

func makeParamsWithArrivalRadius(arrivalCoord util.Coord, arrivalRadius float32, driverOrPassenger string) api.GetJourneysParams {
	params := api.NewGetDriverJourneysParams(util.CoordIgnore, arrivalCoord, 0)
	params.ArrivalRadius = &arrivalRadius

	if driverOrPassenger == "passenger" {
		return castDriverToPassenger(params)
	}

	return params
}

func makeParamsWithTimeDelta(date int, driverOrPassenger string) api.GetJourneysParams {
	params := &api.GetDriverJourneysParams{}
	params.TimeDelta = &date

	if driverOrPassenger == "passenger" {
		return castDriverToPassenger(params)
	}

	return params
}

func makeParamsWithCount(count int, driverOrPassenger string) api.GetJourneysParams {
	params := &api.GetDriverJourneysParams{}
	params.Count = &count

	if driverOrPassenger == "passenger" {
		return castDriverToPassenger(params)
	}

	return params
}
