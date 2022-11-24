package cmd

import (
	"net/http"
	"net/url"
	"time"

	"github.com/fabmob/playground-standard-covoiturage/cmd/test"
	"github.com/spf13/cobra"
)

var postMessagesCmd = makeEndpointCommand(test.PostMessagesEndpoint)

func init() {
	postMessagesCmd.PreRunE = checkGetBookingsCmdFlags

	postMessagesCmd.Run = func(cmd *cobra.Command, args []string) {
		var timeout = 100 * time.Millisecond

		body, err := readBodyFromStdin(cmd, timeout)
		exitWithError(err)

		err = getMessagesRun(
			test.NewDefaultRunner(),
			server,
			body,
		)
		exitWithError(err)
	}

	postCmd.AddCommand(postBookingsCmd)
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
		apiKey,
		flagsWithDefault(http.StatusCreated),
	)
}
