package cmd

import (
	"net/http"
	"net/url"

	"github.com/fabmob/playground-standard-covoiturage/cmd/test"
	"github.com/spf13/cobra"
)

var statusCmd = makeEndpointCommand(test.GetStatusEndpoint)

func init() {

	statusCmd.Run = func(cmd *cobra.Command, args []string) {
		URL, _ := url.JoinPath(server, "/status")
		err := test.RunTest(http.MethodGet, URL, verbose, test.NewQuery(), nil, apiKey, flags(http.StatusOK))
		exitWithError(err)
	}

	getCmd.AddCommand(statusCmd)
}
