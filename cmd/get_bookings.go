package cmd

import (
	"net/http"
	"net/url"

	"github.com/fabmob/playground-standard-covoiturage/cmd/test"
	"github.com/fabmob/playground-standard-covoiturage/cmd/test/endpoint"
	"github.com/spf13/cobra"
)

var (
	getBookingID        string
	expectBookingStatus string
)

var getBookingsCmd = makeEndpointCommand(endpoint.GetBookings)

func init() {

	getBookingsCmd.PreRunE = checkGetBookingsCmdFlags

	getBookingsCmd.Run = func(cmd *cobra.Command, args []string) {
		URL, _ := url.JoinPath(server, "/bookings", getBookingID)
		err := test.RunTest(http.MethodGet, URL, verbose, test.NewQuery(), nil,
			apiKey, flagsWithDefault(http.StatusOK))
		exitWithError(err)
	}

	getBookingsCmd.Flags().StringVar(
		&getBookingID, "bookingId", "", "bookingId path parameter",
	)
	getBookingsCmd.Flags().StringVar(
		&expectBookingStatus, "expectBookingStatus", "", "Expected booking status, checked on response",
	)

	getCmd.AddCommand(getBookingsCmd)
}

func checkGetBookingsCmdFlags(cmd *cobra.Command, args []string) error {
	return checkRequiredBookingID(getBookingID)
}
