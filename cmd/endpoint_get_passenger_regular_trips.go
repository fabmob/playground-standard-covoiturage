package cmd

import (
	"net/http"

	"github.com/fabmob/playground-standard-covoiturage/cmd/endpoint"
)

var getPassengerRegularTripsParameters = getDriverRegularTripsParameters

var _ = makeEndpointCommand(
	endpoint.GetPassengerRegularTrips,
	getPassengerRegularTripsParameters,
	false,
	getCmd,
	http.StatusOK,
)
