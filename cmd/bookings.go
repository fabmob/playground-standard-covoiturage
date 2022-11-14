package cmd

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/fabmob/playground-standard-covoiturage/cmd/test"
	"github.com/spf13/cobra"
)

var (
	bookingID string
	status    string
	message   string
)

//////////////////////////////////////////////////////////////
// GET /bookings/{bookingId}
//////////////////////////////////////////////////////////////

// bookingsCmd represents the bookings command
var getBookingsCmd = &cobra.Command{
	Use:     "bookings",
	Short:   "Test the GET /bookings/{bookingID} endpoint",
	Long:    `Test the GET /bookings/{bookingID} endpoint`,
	PreRunE: checkGetBookingsCmdFlags,
	Run: func(cmd *cobra.Command, args []string) {
		URL, _ := url.JoinPath(server, "/bookings", bookingID)
		err := test.RunTest(http.MethodGet, URL, verbose, test.NewQuery(), nil, flags(http.StatusOK))
		exitWithError(err)
	},
}

func initGetBookings() {
	getBookingsCmd.Flags().StringVar(
		&bookingID, "bookingId", "", "bookingId path parameter",
	)
	getCmd.AddCommand(getBookingsCmd)
}

func checkGetBookingsCmdFlags(cmd *cobra.Command, args []string) error {
	return checkRequiredBookingID(bookingID)
}

//////////////////////////////////////////////////////////////
// POST /bookings
//////////////////////////////////////////////////////////////

// bookingsCmd represents the bookings command
var postBookingsCmd = &cobra.Command{
	Use:   "bookings",
	Short: "Test the POST /bookings endpoint",
	Long:  `Test the POST /bookings endpoint`,
	Run: func(cmd *cobra.Command, args []string) {

		var timeout = 100 * time.Millisecond

		body, err := readBodyFromStdin(cmd, timeout)
		exitWithError(err)

		URL, _ := url.JoinPath(server, "/bookings")
		err = test.RunTest(http.MethodPost, URL, verbose, test.NewQuery(), body, flags(http.StatusCreated))
		exitWithError(err)

	},
}

func initPostBookings() {
	postCmd.AddCommand(postBookingsCmd)
}

//////////////////////////////////////////////////////////////
// PATCH /bookings/{bookingId}
//////////////////////////////////////////////////////////////

// bookingsCmd represents the bookings command
var patchBookingsCmd = &cobra.Command{
	Use:     "bookings",
	Short:   "Test the PATCH /bookings/{bookingID} endpoint",
	Long:    `Test the PATCH /bookings/{bookingID} endpoint`,
	PreRunE: checkPatchBookingsCmdFlags,
	Run: func(cmd *cobra.Command, args []string) {
		err := patchBookingsRun(
			test.NewDefaultRunner(),
			server,
			bookingID,
			status,
			message,
			flags(http.StatusCreated),
		)
		exitWithError(err)
	},
}

func patchBookingsRun(runner test.TestRunner, server, bookingID, status, message string, flags test.Flags) error {

	query := test.NewQuery()
	query.Params["status"] = status
	query.Params["message"] = message

	URL, err := url.JoinPath(server, "/bookings", bookingID)
	if err != nil {
		return err
	}
	return runner.Run(http.MethodPatch, URL, verbose, query, nil, flags)
}

func initPatchBookings() {
	patchBookingsCmd.Flags().StringVar(
		&bookingID, "bookingId", "", "bookingId path parameter",
	)
	patchBookingsCmd.Flags().StringVar(
		&status, "status", "", "status query parameter",
	)
	patchBookingsCmd.Flags().StringVar(
		&message, "message", "", "message query parameter",
	)
	patchCmd.AddCommand(patchBookingsCmd)
}

func checkPatchBookingsCmdFlags(cmd *cobra.Command, args []string) error {
	errRequiredBookingID := checkRequiredBookingID(bookingID)
	if errRequiredBookingID != nil {
		return errRequiredBookingID
	}

	errRequiredStatus := checkRequiredStatus(status)
	if errRequiredStatus != nil {
		return errRequiredStatus
	}

	return nil
}

//////////////////////////////////////////////////////////////
// Other related stuff
//////////////////////////////////////////////////////////////

func init() {
	initGetBookings()
	initPostBookings()
	initPatchBookings()
}

func checkRequired(description string) func(string) error {
	return func(obj string) error {
		if obj == "" {
			return fmt.Errorf("missing required --%s information", description)
		}

		return nil
	}
}

var (
	checkRequiredBookingID = checkRequired("bookingId")
	checkRequiredStatus    = checkRequired("status")
)

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
