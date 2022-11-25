package cmd

import (
	"net/http"
	"net/url"

	"github.com/fabmob/playground-standard-covoiturage/cmd/test"
	"github.com/fabmob/playground-standard-covoiturage/cmd/test/endpoint"
	"github.com/spf13/cobra"
)

var statusCmd = makeEndpointCommand(endpoint.GetStatus)

func init() {

	statusCmd.Run = func(cmd *cobra.Command, args []string) {
		URL, _ := url.JoinPath(server, "/status")
		err := test.RunTest(http.MethodGet, URL, verbose, test.NewQuery(), nil, apiKey, flagsWithDefault(http.StatusOK))
		exitWithError(err)
	}

	getCmd.AddCommand(statusCmd)
}
