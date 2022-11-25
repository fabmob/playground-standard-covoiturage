package cmd

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/fabmob/playground-standard-covoiturage/cmd/test"
	"github.com/spf13/cobra"
)

var driverRegularTripsCmd = makeEndpointCommand(test.GetDriverRegularTripsEndpoint)

var (
	departureTimeOfDay string
	departureWeekdays  []string
	minDepartureDate   string
	maxDepartureDate   string
)

func init() {
	driverRegularTripsCmd.PreRunE = checkGetRegularTripsCmdFlags

	driverRegularTripsCmd.Run = func(cmd *cobra.Command, args []string) {
		err := getDriverRegularTripsRun(
			test.NewDefaultRunner(),
			server,
			departureLat, departureLng, arrivalLat, arrivalLng, departureTimeOfDay,
			timeDelta, departureRadius, arrivalRadius, count,
		)
		exitWithError(err)
	}

	driverRegularTripsCmd.Flags().StringVar(
		&departureLat, "departureLat", "", "departureLat query query parameter")
	driverRegularTripsCmd.Flags().StringVar(
		&departureLng, "departureLng", "", "departureLng query parameter")
	driverRegularTripsCmd.Flags().StringVar(
		&arrivalLat, "arrivalLat", "", "arrivalLat query parameter")
	driverRegularTripsCmd.Flags().StringVar(
		&arrivalLng, "arrivalLng", "", "arrivalLng query parameter")
	driverRegularTripsCmd.Flags().StringVar(
		&departureTimeOfDay, "departureTimeOfDay", "", "departureTimeOfDay query parameter")

	driverRegularTripsCmd.Flags().StringSliceVar(
		&departureWeekdays, "departureWeekdays", []string{}, "departureWeekdays query parameter")
	driverRegularTripsCmd.Flags().StringVar(
		&timeDelta, "timeDelta", "", "timeDelta query parameter")
	driverRegularTripsCmd.Flags().StringVar(
		&departureRadius, "departureRadius", "", "departureRadius query parameter")
	driverRegularTripsCmd.Flags().StringVar(
		&arrivalRadius, "arrivalRadius", "", "arrivalRadius query parameter")
	driverRegularTripsCmd.Flags().StringVar(
		&minDepartureDate, "minDepartureDate", "", "minDepartureDate query parameter")
	driverRegularTripsCmd.Flags().StringVar(
		&maxDepartureDate, "maxDepartureDate", "", "maxDepartureDate query parameter")
	driverRegularTripsCmd.Flags().StringVar(
		&count, "count", "", "count query parameter")

	getCmd.AddCommand(driverRegularTripsCmd)
}

func getDriverRegularTripsRun(
	runner test.TestRunner,
	server string,
	departureLat, departureLng, arrivalLat, arrivalLng, departureTimeOfDay,
	timeDelta, departureRadius, arrivalRadius, count string,
) error {
	query := makeRegularTripQuery(
		departureLat, departureLng, arrivalLat, arrivalLng, departureTimeOfDay,
		departureWeekdays, timeDelta, departureRadius, arrivalRadius, count,
		minDepartureDate, maxDepartureDate,
	)
	URL, _ := url.JoinPath(server, "/driver_regular_trips")

	return runner.Run(http.MethodGet, URL, verbose, query, nil, apiKey, flagsWithDefault(http.StatusOK))
}

func makeRegularTripQuery(
	departureLat, departureLng, arrivalLat, arrivalLng, departureTimeOfDay string,
	departureWeekdays []string,
	timeDelta, departureRadius, arrivalRadius, count, minDepartureDate,
	maxDepartureDate string,
) test.Query {

	var query = test.NewQuery()

	query.SetParam("departureLat", departureLat)
	query.SetParam("departureLng", departureLng)
	query.SetParam("arrivalLat", arrivalLat)
	query.SetParam("arrivalLng", arrivalLng)
	query.SetParam("departureTimeOfDay", departureTimeOfDay)

	departureWeekdaysBytes, _ := json.Marshal(departureWeekdays)
	query.SetOptionalParam("departureWeekdays", string(departureWeekdaysBytes))
	query.SetOptionalParam("timeDelta", timeDelta)
	query.SetOptionalParam("departureRadius", departureRadius)
	query.SetOptionalParam("arrivalRadius", arrivalRadius)
	query.SetOptionalParam("count", count)
	query.SetOptionalParam("minDepartureDate", minDepartureDate)
	query.SetOptionalParam("maxDepartureDate", maxDepartureDate)

	return query
}

func checkGetRegularTripsCmdFlags(cmd *cobra.Command, args []string) error {
	return anyError(
		checkRequiredDepartureLat(departureLat),
		checkRequiredDepartureLng(departureLng),
		checkRequiredArrivalLat(departureLat),
		checkRequiredArrivalLng(departureLng),
		checkRequiredDepartureTimeOfDay(departureTimeOfDay),
	)
}
