package cmd

import (
	"errors"
	"io"
	"net/http"
	"net/url"
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

		var timeout = 100 * time.Millisecond

		body, err := readBodyFromStdin(cmd, timeout)
		exitWithError(err)

		URL, _ := url.JoinPath(server, "/bookings")
		err = test.Run(http.MethodPost, URL, verbose, test.NewQuery(), body, flags(http.StatusCreated))
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

// readBodyFromStdin reads stdin stream until it is closed, and returns its
// content. The function returns an error if it is not closed before `timeout`, or if an error occurs while
// reading.
func readBodyFromStdin(cmd *cobra.Command, timeout time.Duration) ([]byte,
	error) {

	var (
		stdinStreamReader = cmd.InOrStdin()
		stdinChannel      = make(chan []byte, 1)
		errChannel        = make(chan error, 1)
	)

	go func() {
		b, err := io.ReadAll(stdinStreamReader)
		if err != nil {
			errChannel <- err
		} else {
			stdinChannel <- b
		}
	}()

	select {
	case <-time.After(timeout):
		return nil, errors.New("body is required but missing")

	case err := <-errChannel:
		return nil, err

	case body := <-stdinChannel:
		return body, nil
	}
}
