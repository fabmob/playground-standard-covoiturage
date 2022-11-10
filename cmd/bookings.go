package cmd

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

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
		err := test.Run(http.MethodGet, URL, verbose, test.NewQuery(), nil, flags(http.StatusOK))
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
	Use:   "bookings",
	Short: "Test the GET /bookings endpoint",
	Long:  `Test the GET /bookings endpoint`,
	Run: func(cmd *cobra.Command, args []string) {

		body := cmd.InOrStdin()
		stdinChannel := make(chan []byte, 1)

		go func() {
			b, _ := io.ReadAll(body)
			stdinChannel <- b
		}()

		var timeout = 100 * time.Millisecond

		select {
		case <-time.After(timeout):
			fmt.Println("body is required but missing")
			os.Exit(1)

		case body := <-stdinChannel:
			URL, _ := url.JoinPath(server, "/bookings")
			fmt.Println(URL)
			err := test.Run(http.MethodPost, URL, verbose, test.NewQuery(), body, flags(http.StatusCreated))
			exitWithError(err)
		}
	},
}

func initPostBookings() {
	postCmd.AddCommand(postBookingsCmd)
}

func init() {
	initGetBookings()
	initPostBookings()
}
