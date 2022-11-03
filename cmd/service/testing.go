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

func makeParamsWithDepartureRadius(departureCoord util.Coord, departureRadius float32) *api.GetDriverJourneysParams {
	params := api.NewGetDriverJourneysParams(departureCoord, util.CoordIgnore, 0)
	params.DepartureRadius = &departureRadius
	return params
}

func makeParamsWithDepartureRadius2(departureCoord util.Coord, departureRadius float32) *api.GetPassengerJourneysParams {
	params := api.NewGetPassengerJourneysParams(departureCoord, util.CoordIgnore, 0)
	params.DepartureRadius = &departureRadius
	return params
}

func makeParamsWithArrivalRadius(arrivalCoord util.Coord, arrivalRadius float32) *api.GetDriverJourneysParams {
	params := api.NewGetDriverJourneysParams(util.CoordIgnore, arrivalCoord, 0)
	params.ArrivalRadius = &arrivalRadius
	return params
}

func makeParamsWithTimeDelta(date int) *api.GetDriverJourneysParams {
	params := &api.GetDriverJourneysParams{}
	params.TimeDelta = &date
	return params
}

func makeParamsWithCount(count int) *api.GetDriverJourneysParams {
	params := &api.GetDriverJourneysParams{}
	params.Count = &count
	return params
}
