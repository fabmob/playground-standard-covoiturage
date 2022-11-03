package service

import (
	"gitlab.com/multi/stdcov-api-test/cmd/api"
	"gitlab.com/multi/stdcov-api-test/cmd/util"
)

func makeNDriverJourneys(n int) []api.DriverJourney {
	driverJourneys := make([]api.DriverJourney, 0, n)
	for i := 0; i < n; i++ {
		driverJourneys = append(driverJourneys, api.NewDriverJourney())
	}
	return driverJourneys
}

func makeDriverJourneyAtCoords(coordPickup, coordDrop util.Coord) api.DriverJourney {
	dj := api.NewDriverJourney()
	dj.PassengerPickupLat = coordPickup.Lat
	dj.PassengerPickupLng = coordPickup.Lon
	dj.PassengerDropLat = coordDrop.Lat
	dj.PassengerDropLng = coordDrop.Lon
	return dj
}

func makePassengerJourneyAtCoords(coordPickup, coordDrop util.Coord) api.PassengerJourney {
	pj := api.NewPassengerJourney()
	pj.PassengerPickupLat = coordPickup.Lat
	pj.PassengerPickupLng = coordPickup.Lon
	pj.PassengerDropLat = coordDrop.Lat
	pj.PassengerDropLng = coordDrop.Lon
	return pj
}

func makeDriverJourneyAtDate(date int64) api.DriverJourney {
	dj := api.NewDriverJourney()
	dj.PassengerPickupDate = date
	return dj
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
