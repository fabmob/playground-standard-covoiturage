package cmd

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/fabmob/playground-standard-covoiturage/cmd/endpoint"
	"github.com/fabmob/playground-standard-covoiturage/cmd/test"
	"github.com/spf13/cobra"
)

var driverRegularTripsCmd = makeEndpointCommand(endpoint.GetDriverRegularTrips)

var (
	departureTimeOfDay string
	departureWeekdays  []string
	minDepartureDate   string
	maxDepartureDate   string
)

// only string parameters. Array parameter departureWeekdays is treated
// separately.
var getDriverRegularTripsParameters = []parameter{
	{&departureLat, "departureLat", true, "query"},
	{&departureLng, "departureLng", true, "query"},
	{&arrivalLat, "arrivalLat", true, "query"},
	{&arrivalLng, "arrivalLng", true, "query"},
	{&departureTimeOfDay, "departureTimeOfDay", true, "query"},
	{&timeDelta, "timeDelta", false, "query"},
	{&departureRadius, "departureRadius", false, "query"},
	{&arrivalRadius, "arrivalRadius", false, "query"},
	{&minDepartureDate, "minDepartureDate", false, "query"},
	{&maxDepartureDate, "maxDepartureDate", false, "query"},
	{&count, "count", false, "query"},
}

func init() {
	cmd := driverRegularTripsCmd
	cmd.PreRunE = checkGetRegularTripsCmdFlags

	cmd.Run = func(cmd *cobra.Command, args []string) {
		err := getRegularTripsRun(
			test.NewDefaultRunner(),
			server,
			getDriverRegularTripsParameters,
			departureWeekdays,
			"/driver_regular_trips",
		)
		exitWithError(err)
	}

	for _, q := range getDriverRegularTripsParameters {
		parameterFlag(cmd.Flags(), q.where, q.variable, q.name, q.required)
	}

	// additional array query parameter
	cmd.Flags().StringSliceVar(
		&departureWeekdays, "departureWeekdays", []string{}, "departureWeekdays query parameter")

	getCmd.AddCommand(cmd)
}

func getRegularTripsRun(
	runner test.TestRunner,
	server string,
	queryStringParameters []parameter,
	departureWeekdays []string,
	endpointPath string,
) error {
	query := makeQuery(queryStringParameters)

	departureWeekdaysBytes, err := json.Marshal(departureWeekdays)
	if err != nil {
		return err
	}
	query.SetOptionalParam("departureWeekdays", string(departureWeekdaysBytes))

	URL, err := url.JoinPath(server, endpointPath)
	if err != nil {
		return err
	}

	return runner.Run(http.MethodGet, URL, query, nil, verbose, apiKey, flagsWithDefault(http.StatusOK))
}

func checkGetRegularTripsCmdFlags(cmd *cobra.Command, args []string) error {
	for _, q := range getDriverRegularTripsParameters {
		if err := checkRequired(q.variable, q.name); err != nil {
			return err
		}
	}

	return nil
}
