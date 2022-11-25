package cmd

import (
	"net/http"
	"net/url"

	"github.com/fabmob/playground-standard-covoiturage/cmd/endpoint"
	"github.com/fabmob/playground-standard-covoiturage/cmd/test"
	"github.com/spf13/cobra"
)

var driverJourneysCmd = makeEndpointCommand(endpoint.GetDriverJourneys)

var (
	departureLat    string
	departureLng    string
	arrivalLat      string
	arrivalLng      string
	departureDate   string
	timeDelta       string
	departureRadius string
	arrivalRadius   string
	count           string
)

var getDriverJourneysParameters = []parameter{
	{&departureLat, "departureLat", true, "query"},
	{&departureLng, "departureLng", true, "query"},
	{&arrivalLat, "arrivalLat", true, "query"},
	{&arrivalLng, "arrivalLng", true, "query"},
	{&departureDate, "departureDate", true, "query"},
	{&timeDelta, "timeDelta", false, "query"},
	{&departureRadius, "departureRadius", false, "query"},
	{&arrivalRadius, "arrivalRadius", false, "query"},
	{&count, "count", false, "query"},
}

func init() {
	cmd := driverJourneysCmd
	cmd.PreRunE = checkGetJourneysCmdFlags

	cmd.Run = func(cmd *cobra.Command, args []string) {
		query := makeQuery(getDriverJourneysParameters)
		URL, _ := url.JoinPath(server, "/driver_journeys")
		err := test.RunTest(http.MethodGet, URL, verbose, query, nil, apiKey, flagsWithDefault(http.StatusOK))
		exitWithError(err)
	}

	for _, q := range getDriverJourneysParameters {
		parameterFlag(cmd.Flags(), q.where, q.variable, q.name, q.required)
	}

	getCmd.AddCommand(cmd)
}

func checkGetJourneysCmdFlags(cmd *cobra.Command, args []string) error {
	for _, q := range getDriverJourneysParameters {
		if err := checkRequired(q.variable, q.name); err != nil {
			return err
		}
	}

	return nil
}
