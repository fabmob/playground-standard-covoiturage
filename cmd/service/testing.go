package service

import (
	"gitlab.com/multi/stdcov-api-test/cmd/api"
	"gitlab.com/multi/stdcov-api-test/cmd/util"
)

func makeDriverJourney(coordPickup, coordDrop util.Coord) api.DriverJourney {
	return api.DriverJourney{
		PassengerPickupLat: coordPickup.Lat,
		PassengerPickupLng: coordPickup.Lon,
		PassengerDropLat:   coordDrop.Lat,
		PassengerDropLng:   coordDrop.Lon,
		Type:               "DYNAMIC",
	}
}

func makeParamsWithDepartureRadius(departureCoord util.Coord, departureRadius float32) *api.GetDriverJourneysParams {
	params := api.NewGetDriverJourneysParams(
		float32(departureCoord.Lat),
		float32(departureCoord.Lon),
		0,
		0,
		0,
	)
	params.DepartureRadius = &departureRadius
	return params
}

func makeParamsWithArrivalRadius(arrivalCoord util.Coord, arrivalRadius float32) *api.GetDriverJourneysParams {
	params := api.NewGetDriverJourneysParams(
		float32(arrivalCoord.Lat),
		float32(arrivalCoord.Lon),
		0,
		0,
		0,
	)
	params.ArrivalRadius = &arrivalRadius
	return params
}
