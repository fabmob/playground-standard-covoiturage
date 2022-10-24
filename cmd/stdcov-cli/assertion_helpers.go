package main

import (
	"fmt"

	"github.com/umahmood/haversine"
	"gitlab.com/multi/stdcov-api-test/cmd/stdcov-cli/client"
)

type coords struct {
	lat float64
	lon float64
}

func distanceKm(coords1, coords2 coords) float64 {
	c1 := haversine.Coord{Lat: coords1.lat, Lon: coords1.lon}
	c2 := haversine.Coord{Lat: coords2.lat, Lon: coords2.lon}
	_, dist := haversine.Distance(c1, c2)
	return dist
}

// getQueryCoords extracts departure or arrival coordinates from
// queryParameters
func getQueryCoords(departureOrArrival departureOrArrival, queryParams *client.GetDriverJourneysParams) coords {
	var coordsQuery coords
	switch departureOrArrival {
	case departure:
		coordsQuery = coords{float64(queryParams.DepartureLat), float64(queryParams.DepartureLng)}
	case arrival:
		coordsQuery = coords{float64(queryParams.ArrivalLat), float64(queryParams.ArrivalLng)}
	}
	return coordsQuery
}

// getResponseCoords extracts departure or arrival coordinates from
// driverJourney object. Fails if required coordinates are missing.
func getResponseCoords(departureOrArrival departureOrArrival, driverJourney client.DriverJourney) (coords, error) {
	missingDeparture := departureOrArrival == departure &&
		(driverJourney.DriverDepartureLng == nil || driverJourney.DriverDepartureLat == nil)
	missingArrival := departureOrArrival == arrival &&
		(driverJourney.DriverArrivalLng == nil || driverJourney.DriverArrivalLat ==
			nil)
	if missingDeparture || missingArrival {

		return coords{}, fmt.Errorf("malformed response: driverDepartureLat, driverDepartureLng, driverArrivalLat and driverArrivalLng are required")
	}
	var coordsResponse coords
	switch departureOrArrival {
	case departure:
		coordsResponse = coords{*driverJourney.DriverDepartureLat, *driverJourney.DriverDepartureLng}
	case arrival:
		coordsResponse = coords{*driverJourney.DriverArrivalLat, *driverJourney.DriverArrivalLng}
	}
	return coordsResponse, nil
}

// getQueryRadiusOrDefault returns departureRadius er arrivalRadius query parameter
// (depending on departureOrArrival), or the default value if missing
func getQueryRadiusOrDefault(departureOrArrival departureOrArrival, queryParams *client.GetDriverJourneysParams) float64 {
	const DefaultRadius float32 = 1
	var radiusPtr *float32
	switch departureOrArrival {
	case departure:
		radiusPtr = queryParams.DepartureRadius
	case arrival:
		radiusPtr = queryParams.ArrivalRadius
	}
	var radius float32
	if radiusPtr != nil {
		radius = *radiusPtr
	} else {
		radius = DefaultRadius
	}
	return float64(radius)
}

// failedParsing wraps a parsing error with additional details
func failedParsing(responseOrRequest string, err error) error {
	return fmt.Errorf(
		"internal error while parsing %s:%w",
		responseOrRequest,
		err,
	)
}
