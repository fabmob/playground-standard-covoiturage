package cmd

import (
	"net/http"

	"github.com/fabmob/playground-standard-covoiturage/cmd/endpoint"
)

var (
	departureLat    string
	departureLng    string
	arrivalLat      string
	arrivalLng      string
	departureDate   string
	timeDelta       string
	departureRadius string
	arrivalRadius   string
	count           string
)

var getDriverJourneysParameters = parameters{
	{&departureLat, "departureLat", true, "query"},
	{&departureLng, "departureLng", true, "query"},
	{&arrivalLat, "arrivalLat", true, "query"},
	{&arrivalLng, "arrivalLng", true, "query"},
	{&departureDate, "departureDate", true, "query"},
	{&timeDelta, "timeDelta", false, "query"},
	{&departureRadius, "departureRadius", false, "query"},
	{&arrivalRadius, "arrivalRadius", false, "query"},
	{&count, "count", false, "query"},
}

var _ = makeEndpointCommand(
	endpoint.GetDriverJourneys,
	getDriverJourneysParameters,
	false,
	getCmd,
	http.StatusOK,
)
