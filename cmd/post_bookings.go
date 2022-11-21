package cmd

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/fabmob/playground-standard-covoiturage/cmd/test"
	"github.com/spf13/cobra"
)

var postBookingsCmd = makeEndpointCommand(test.PostBookingsEndpoint)

func init() {
	postBookingsCmd.Run = func(cmd *cobra.Command, args []string) {
		var timeout = 100 * time.Millisecond

		body, err := readBodyFromStdin(cmd, timeout)
		exitWithError(err)

		fmt.Println(string(body))

		URL, _ := url.JoinPath(server, "/bookings")
		err = test.RunTest(http.MethodPost, URL, verbose, test.NewQuery(), body, apiKey, flags(http.StatusCreated))
		exitWithError(err)
	}

	postCmd.AddCommand(postBookingsCmd)
}
