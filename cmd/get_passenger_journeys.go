package cmd

import (
	"net/http"
	"net/url"

	"github.com/fabmob/playground-standard-covoiturage/cmd/endpoint"
	"github.com/fabmob/playground-standard-covoiturage/cmd/test"
	"github.com/spf13/cobra"
)

var passengerJourneysCmd = makeEndpointCommand(endpoint.GetPassengerJourneys)

var getPassengerJourneysParameters = getDriverJourneysParameters

func init() {
	cmd := passengerJourneysCmd
	cmd.PreRunE = checkRequiredCmdFlags(getPassengerJourneysParameters)

	cmd.Run = func(cmd *cobra.Command, args []string) {
		query := makeQuery(getPassengerJourneysParameters)
		URL, _ := url.JoinPath(server, "/passenger_journeys")
		err := test.RunTest(http.MethodGet, URL, verbose, query, nil, apiKey, flagsWithDefault(http.StatusOK))
		exitWithError(err)
	}

	for _, q := range getDriverJourneysParameters {
		parameterFlag(cmd.Flags(), q.where, q.variable, q.name, q.required)
	}

	getCmd.AddCommand(cmd)
}
