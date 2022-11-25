package cmd

import (
	"net/http"

	"github.com/fabmob/playground-standard-covoiturage/cmd/endpoint"
)

var (
	getBookingID          string
	expectedBookingStatus string
)

var getBookingsParameters = parameters{
	{&getBookingID, "bookingId", true, "path"},
}

var getBookingsCmd = makeEndpointCommand(
	endpoint.GetBookings,
	getBookingsParameters,
	false,
	getCmd,
	http.StatusOK,
)

func init() {
	cmd := getBookingsCmd
	cmd.Flags().StringVar(
		&expectedBookingStatus, "expectedBookingStatus", "", "Expected booking status, checked on response",
	)
}
