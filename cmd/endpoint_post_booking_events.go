package cmd

import (
	"net/http"

	"github.com/fabmob/playground-standard-covoiturage/cmd/endpoint"
)

var postBookingEventsParameters = parameters{}

var _ = makeEndpointCommand(
	endpoint.PostBookingEvents,
	postBookingEventsParameters,
	true,
	postCmd,
	http.StatusOK,
)
