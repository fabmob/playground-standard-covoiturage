package cmd

import (
	"net/http"

	"github.com/fabmob/playground-standard-covoiturage/cmd/endpoint"
)

var (
	departureTimeOfDay   string
	departureWeekdaysStr string
	minDepartureDate     string
	maxDepartureDate     string
)

var getDriverRegularTripsParameters = parameters{
	{&departureLat, "departureLat", true, "query"},
	{&departureLng, "departureLng", true, "query"},
	{&arrivalLat, "arrivalLat", true, "query"},
	{&arrivalLng, "arrivalLng", true, "query"},
	{&departureTimeOfDay, "departureTimeOfDay", true, "query"},
	{&timeDelta, "timeDelta", false, "query"},
	{&departureRadius, "departureRadius", false, "query"},
	{&arrivalRadius, "arrivalRadius", false, "query"},
	{&minDepartureDate, "minDepartureDate", false, "query"},
	{&maxDepartureDate, "maxDepartureDate", false, "query"},
	{&count, "count", false, "query"},
	{&departureWeekdaysStr, "departureWeekdays", false, "query"},
}

var driverRegularTripsCmd = makeEndpointCommand(
	endpoint.GetDriverRegularTrips,
	getDriverRegularTripsParameters,
	false,
	getCmd,
	http.StatusOK,
)
