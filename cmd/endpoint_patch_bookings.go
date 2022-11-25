package cmd

import (
	"net/http"

	"github.com/fabmob/playground-standard-covoiturage/cmd/endpoint"
)

var (
	patchBookingID string
	status         string
	message        string
)

var patchBookingsParameters = []parameter{
	{&status, "status", true, "query"},
	{&message, "message", false, "query"},
	{&patchBookingID, "bookingId", true, "path"},
}

var _ = makeEndpointCommand(
	endpoint.PatchBookings,
	patchBookingsParameters,
	false,
	patchCmd,
	http.StatusOK,
)
