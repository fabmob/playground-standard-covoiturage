package cmd

import (
	"net/http"
	"net/url"

	"github.com/fabmob/playground-standard-covoiturage/cmd/test"
	"github.com/spf13/cobra"
)

//////////////////////////////////////////////////////////////
// GET /bookings/{bookingId}
//////////////////////////////////////////////////////////////

var getBookingID string

// bookingsCmd represents the bookings command
var getBookingsCmd = &cobra.Command{
	Use:     "bookings",
	Short:   "Test the GET /bookings/{bookingID} endpoint",
	Long:    `Test the GET /bookings/{bookingID} endpoint`,
	PreRunE: checkGetBookingsCmdFlags,
	Run: func(cmd *cobra.Command, args []string) {
		URL, _ := url.JoinPath(server, "/bookings", getBookingID)
		err := test.RunTest(http.MethodGet, URL, verbose, test.NewQuery(), nil, flags(http.StatusOK))
		exitWithError(err)
	},
}

func init() {
	getBookingsCmd.Flags().StringVar(
		&getBookingID, "bookingId", "", "bookingId path parameter",
	)
	getCmd.AddCommand(getBookingsCmd)
}

func checkGetBookingsCmdFlags(cmd *cobra.Command, args []string) error {
	return checkRequiredBookingID(getBookingID)
}
