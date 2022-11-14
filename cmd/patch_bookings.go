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

var patchBookingsCmd = &cobra.Command{
	Use:     "bookings",
	Short:   cmdDescription(test.PatchBookingsEndpoint),
	Long:    cmdDescription(test.PatchBookingsEndpoint),
	PreRunE: checkPatchBookingsCmdFlags,
	Run: func(cmd *cobra.Command, args []string) {
		err := patchBookingsRun(
			test.NewDefaultRunner(),
			server,
			patchBookingID,
			status,
			message,
		)
		exitWithError(err)
	},
}

func patchBookingsRun(runner test.TestRunner, server, bookingID, status, message string) error {

	query := test.NewQuery()
	query.Params["status"] = status
	query.Params["message"] = message

	URL, err := url.JoinPath(server, "/bookings", bookingID)
	if err != nil {
		return err
	}

	return runner.Run(http.MethodPatch, URL, verbose, query, nil, flags(http.StatusOK))
}

func init() {
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
