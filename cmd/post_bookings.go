package cmd

import (
	"net/http"
	"net/url"
	"time"

	"github.com/fabmob/playground-standard-covoiturage/cmd/test"
	"github.com/fabmob/playground-standard-covoiturage/cmd/test/endpoint"
	"github.com/spf13/cobra"
)

var postBookingsCmd = makeEndpointCommand(endpoint.PostBookings)

func init() {
	postBookingsCmd.Run = func(cmd *cobra.Command, args []string) {
		var timeout = 100 * time.Millisecond

		body, err := readBodyFromStdin(cmd, timeout)
		exitWithError(err)

		URL, _ := url.JoinPath(server, "/bookings")
		err = test.RunTest(http.MethodPost, URL, verbose, test.NewQuery(), body, apiKey, flagsWithDefault(http.StatusCreated))
		exitWithError(err)
	}

	postCmd.AddCommand(postBookingsCmd)
}
