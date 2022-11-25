package cmd

import (
	"net/http"

	"github.com/fabmob/playground-standard-covoiturage/cmd/endpoint"
)

var postBookingsParameters = []parameter{}

var postBookingsCmd = makeEndpointCommand(
	endpoint.PostBookings,
	postBookingsParameters,
	true,
	postCmd,
	http.StatusCreated,
)
