package cmd

import (
	"net/http"
	"net/url"

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
		err := getMessagesRun(
			test.NewDefaultRunner(),
			server,
		)
		exitWithError(err)
	},
}

func getMessagesRun(runner test.TestRunner, server string) error {
	URL, _ := url.JoinPath(server, "/messages")
	return runner.Run(http.MethodGet, URL, verbose, test.NewQuery(), nil, flags(http.StatusOK))
}

func init() {
	postCmd.AddCommand(postBookingsCmd)
}
