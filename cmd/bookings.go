package cmd

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/fabmob/playground-standard-covoiturage/cmd/test"
	"github.com/spf13/cobra"
)

// bookingsCmd represents the bookings command
var bookingsCmd = &cobra.Command{
	Use:     "bookings",
	Short:   "Test the GET /bookings endpoint",
	Long:    `Test the GET /bookings endpoint`,
	PreRunE: checkGetBookingsCmdFlags,
	Run: func(cmd *cobra.Command, args []string) {
		URL, _ := url.JoinPath(server, "/bookings", bookingID)
		err := test.Run(http.MethodGet, URL, verbose, test.NewQuery(), flags())
		exitWithError(err)
	},
}

var (
	bookingID string
)

func init() {
	bookingsCmd.Flags().StringVar(
		&bookingID, "bookingId", "", "bookingIdpath parameter",
	)
	getCmd.AddCommand(bookingsCmd)
}

func checkGetBookingsCmdFlags(cmd *cobra.Command, args []string) error {
	if bookingID == "" {
		return errors.New("missing required --bookingId information")
	}

	return nil
}
