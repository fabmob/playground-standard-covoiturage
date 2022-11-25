package cmd

import (
	"net/http"
	"net/url"
	"time"

	"github.com/fabmob/playground-standard-covoiturage/cmd/test"
	"github.com/spf13/cobra"
)

var postBookingEventsCmd = makeEndpointCommand(test.PostBookingEventsEndpoint)

func init() {
	postBookingEventsCmd.Run = func(cmd *cobra.Command, args []string) {
		var timeout = 100 * time.Millisecond

		body, err := readBodyFromStdin(cmd, timeout)
		exitWithError(err)

		err = postBookingEventsRun(
			test.NewDefaultRunner(),
			server,
			body,
		)
		exitWithError(err)
	}

	postCmd.AddCommand(postBookingsCmd)
}

func postBookingEventsRun(runner test.TestRunner, server string, body []byte) error {

	URL, err := url.JoinPath(server, "/booking_events")
	if err != nil {
		return err
	}

	return runner.Run(http.MethodPost, URL, verbose, query, body, apiKey, flagsWithDefault(http.StatusOK))
}
