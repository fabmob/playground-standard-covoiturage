package cmd

import (
	"net/http"
	"net/url"

	"github.com/fabmob/playground-standard-covoiturage/cmd/endpoint"
	"github.com/fabmob/playground-standard-covoiturage/cmd/test"
	"github.com/spf13/cobra"
)

var (
	getBookingID          string
	expectedBookingStatus string
)

var getBookingsParameters = []parameter{
	{&getBookingID, "bookingId", true, "path"},
}

var getBookingsCmd = makeEndpointCommand(endpoint.GetBookings)

func init() {
	cmd := getBookingsCmd

	cmd.PreRunE = checkRequiredCmdFlags(getBookingsParameters)

	cmd.Run = func(cmd *cobra.Command, args []string) {
		URL, _ := url.JoinPath(server, "/bookings", getBookingID)
		err := test.RunTest(http.MethodGet, URL, verbose, test.NewQuery(), nil,
			apiKey, flagsWithDefault(http.StatusOK))
		exitWithError(err)
	}

	for _, q := range getBookingsParameters {
		parameterFlag(cmd.Flags(), q.where, q.variable, q.name, q.required)
	}

	cmd.Flags().StringVar(
		&expectedBookingStatus, "expectedBookingStatus", "", "Expected booking status, checked on response",
	)

	getCmd.AddCommand(getBookingsCmd)
}
