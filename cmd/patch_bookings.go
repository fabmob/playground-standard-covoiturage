package cmd

import (
	"net/http"
	"net/url"

	"github.com/fabmob/playground-standard-covoiturage/cmd/test"
	"github.com/spf13/cobra"
)

var (
	patchBookingID string
	status         string
	message        string
)

var patchBookingsCmd = makeEndpointCommand(test.PatchBookingsEndpoint)

func init() {
	patchBookingsCmd.PreRunE = checkPatchBookingsCmdFlags

	patchBookingsCmd.Run = func(cmd *cobra.Command, args []string) {
		err := patchBookingsRun(
			test.NewDefaultRunner(),
			server,
			patchBookingID,
			status,
			message,
		)
		exitWithError(err)
	}

	patchBookingsCmd.Flags().StringVar(
		&patchBookingID, "bookingId", "", "bookingId path parameter",
	)

	patchBookingsCmd.Flags().StringVar(
		&status, "status", "", "status query parameter",
	)

	patchBookingsCmd.Flags().StringVar(
		&message, "message", "", "message query parameter",
	)

	patchCmd.AddCommand(patchBookingsCmd)
}

func patchBookingsRun(runner test.TestRunner, server, bookingID, status, message string) error {
	query := test.NewQuery()
	query.Params["status"] = status
	query.Params["message"] = message

	URL, err := url.JoinPath(server, "/bookings", bookingID)
	if err != nil {
		return err
	}

	return runner.Run(http.MethodPatch, URL, verbose, query, nil, apiKey, flagsWithDefault(http.StatusOK))
}

func checkPatchBookingsCmdFlags(cmd *cobra.Command, args []string) error {
	errRequiredBookingID := checkRequiredBookingID(patchBookingID)
	if errRequiredBookingID != nil {
		return errRequiredBookingID
	}

	errRequiredStatus := checkRequiredStatus(status)
	if errRequiredStatus != nil {
		return errRequiredStatus
	}

	return nil
}
