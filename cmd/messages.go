package cmd

import (
	"net/http"
	"net/url"
	"time"

	"github.com/fabmob/playground-standard-covoiturage/cmd/test"
	"github.com/spf13/cobra"
)

// bookingsCmd represents the bookings command
var getMessagesCmd = &cobra.Command{
	Use:     "bookings",
	Short:   "Test the POST /messages endpoint",
	Long:    `Test the POST /messages endpoint`,
	PreRunE: checkGetBookingsCmdFlags,
	Run: func(cmd *cobra.Command, args []string) {
		var timeout = 100 * time.Millisecond

		body, err := readBodyFromStdin(cmd, timeout)
		exitWithError(err)

		err = getMessagesRun(
			test.NewDefaultRunner(),
			server,
			body,
		)
		exitWithError(err)
	},
}

func getMessagesRun(runner test.TestRunner, server string, body []byte) error {
	URL, _ := url.JoinPath(server, "/messages")

	return runner.Run(http.MethodGet, URL, verbose, test.NewQuery(), body, flags(http.StatusOK))
}

func init() {
	postCmd.AddCommand(postBookingsCmd)
}
