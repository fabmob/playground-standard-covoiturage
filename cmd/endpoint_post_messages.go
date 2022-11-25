package cmd

import (
	"net/http"

	"github.com/fabmob/playground-standard-covoiturage/cmd/endpoint"
)

var postMessagesParameters = []parameter{}

var postMessagesCmd = makeEndpointCommand(
	endpoint.PostMessages,
	postMessagesParameters,
	true,
	postCmd,
	http.StatusOK,
)
