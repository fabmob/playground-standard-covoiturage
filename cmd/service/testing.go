package service

import (
	"gitlab.com/multi/stdcov-api-test/cmd/api"
	"gitlab.com/multi/stdcov-api-test/cmd/util"
)

func makeNDriverJourneys(n int) []api.DriverJourney {
	driverJourneys := make([]api.DriverJourney, 0, n)
	for i := 0; i < n; i++ {
		driverJourneys = append(driverJourneys, api.DriverJourney{Type: "DYNAMIC"})
	}
	return driverJourneys
}

func makeDriverJourneyAtCoords(coordPickup, coordDrop util.Coord) api.DriverJourney {
	return api.DriverJourney{
		PassengerPickupLat: coordPickup.Lat,
		PassengerPickupLng: coordPickup.Lon,
		PassengerDropLat:   coordDrop.Lat,
		PassengerDropLng:   coordDrop.Lon,
		Type:               "DYNAMIC",
	}
}

func makeDriverJourneyAtDate(date int64) api.DriverJourney {
	return api.DriverJourney{
		PassengerPickupDate: date,
		Type:                "DYNAMIC",
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
		0,
		0,
		float32(arrivalCoord.Lat),
		float32(arrivalCoord.Lon),
		0,
	)
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
