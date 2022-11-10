package cmd

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/fabmob/playground-standard-covoiturage/cmd/test"
	"github.com/spf13/cobra"
)

// bookingsCmd represents the bookings command
var getBookingsCmd = &cobra.Command{
	Use:     "bookings",
	Short:   "Test the GET /bookings endpoint",
	Long:    `Test the GET /bookings endpoint`,
	PreRunE: checkGetBookingsCmdFlags,
	Run: func(cmd *cobra.Command, args []string) {
		URL, _ := url.JoinPath(server, "/bookings", bookingID)
		err := test.Run(http.MethodGet, URL, verbose, test.NewQuery(), flags(http.StatusOK))
		exitWithError(err)
	},
}

var (
	bookingID string
)

func initGetBookings() {
	getBookingsCmd.Flags().StringVar(
		&bookingID, "bookingId", "", "bookingId path parameter",
	)
	getCmd.AddCommand(getBookingsCmd)
}

func checkGetBookingsCmdFlags(cmd *cobra.Command, args []string) error {
	if bookingID == "" {
		return errors.New("missing required --bookingId information")
	}

	return nil
}

// bookingsCmd represents the bookings command
var postBookingsCmd = &cobra.Command{
	Use:     "bookings",
	Short:   "Test the GET /bookings endpoint",
	Long:    `Test the GET /bookings endpoint`,
	PreRunE: checkGetBookingsCmdFlags,
	Run: func(cmd *cobra.Command, args []string) {
		url.JoinPath(server, "/bookings")
		err := test.Run(http.MethodPost, URL, verbose, test.NewQuery(), flags(http.StatusCreated))
		exitWithError(err)
	},
}

func initPostBookings() {
	postCmd.AddCommand(postBookingsCmd)
}

func init() {
	initGetBookings()
	initPostBookings()
}
