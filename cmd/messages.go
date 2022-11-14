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
	Short:   cmdDescription(test.PostMessagesEndpoint),
	Long:    cmdDescription(test.PostMessagesEndpoint),
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
	URL, err := url.JoinPath(server, "/messages")
	if err != nil {
		return err
	}
	return runner.Run(
		http.MethodPost,
		URL,
		verbose,
		test.NewQuery(),
		body,
		flags(http.StatusCreated),
	)
}

func init() {
	postCmd.AddCommand(postBookingsCmd)
}
