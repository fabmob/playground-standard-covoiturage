package cmd

import (
	"net/http"
	"net/url"
	"time"

	"github.com/fabmob/playground-standard-covoiturage/cmd/endpoint"
	"github.com/fabmob/playground-standard-covoiturage/cmd/test"
	"github.com/spf13/cobra"
)

var postMessagesCmd = makeEndpointCommand(endpoint.PostMessages)

var postMessagesParameters = []parameter{}

func init() {
	cmd := postMessagesCmd
	cmd.PreRunE = checkRequiredCmdFlags(postMessagesParameters)

	cmd.Run = func(cmd *cobra.Command, args []string) {
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

	postCmd.AddCommand(cmd)
}

func getMessagesRun(runner test.TestRunner, server string, body []byte) error {
	URL, err := url.JoinPath(server, "/messages")
	if err != nil {
		return err
	}
	return runner.Run(
		http.MethodPost,
		URL,
		test.NewQuery(),
		body,
		verbose,
		apiKey,
		flagsWithDefault(http.StatusCreated),
	)
}
