package cmd

import (
	"net/http"

	"github.com/fabmob/playground-standard-covoiturage/cmd/endpoint"
)

var statusCmdParameters = parameters{}

var _ = makeEndpointCommand(
	endpoint.GetStatus,
	statusCmdParameters,
	false,
	getCmd,
	http.StatusOK,
)
