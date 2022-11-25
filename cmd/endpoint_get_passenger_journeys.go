package cmd

import (
	"net/http"

	"github.com/fabmob/playground-standard-covoiturage/cmd/endpoint"
)

var getPassengerJourneysParameters = getDriverJourneysParameters

var _ = makeEndpointCommand(
	endpoint.GetPassengerJourneys,
	getPassengerJourneysParameters,
	false,
	getCmd,
	http.StatusOK,
)
