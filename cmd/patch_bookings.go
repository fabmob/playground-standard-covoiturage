package cmd

import (
	"net/http"
	"net/url"

	"github.com/fabmob/playground-standard-covoiturage/cmd/endpoint"
	"github.com/fabmob/playground-standard-covoiturage/cmd/test"
	"github.com/spf13/cobra"
)

var (
	patchBookingID string
	status         string
	message        string
)

var patchBookingsCmd = makeEndpointCommand(endpoint.PatchBookings)

var patchBookingsParameters = []parameter{
	{&status, "status", true, "query"},
	{&message, "message", false, "query"},
	{&patchBookingID, "bookingId", true, "path"},
}

func init() {
	cmd := patchBookingsCmd
	cmd.PreRunE = checkPatchBookingsCmdFlags

	cmd.Run = func(cmd *cobra.Command, args []string) {
		err := patchBookingsRun(
			test.NewDefaultRunner(),
			server,
			patchBookingID,
			patchBookingsParameters,
		)
		exitWithError(err)
	}

	for _, q := range patchBookingsParameters {
		parameterFlag(cmd.Flags(), q.where, q.variable, q.name, q.required)
	}

	patchCmd.AddCommand(cmd)
}

func patchBookingsRun(runner test.TestRunner, server string, bookingID string, queryParameters []parameter) error {
	query := makeQuery(queryParameters)

	URL, err := url.JoinPath(server, "/bookings", bookingID)
	if err != nil {
		return err
	}

	return runner.Run(http.MethodPatch, URL, query, nil, verbose, apiKey, flagsWithDefault(http.StatusOK))
}

func checkPatchBookingsCmdFlags(cmd *cobra.Command, args []string) error {
	for _, q := range patchBookingsParameters {
		if err := checkRequired(q.variable, q.name); err != nil {
			return err
		}
	}

	return nil
}
