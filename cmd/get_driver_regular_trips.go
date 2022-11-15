package cmd

import (
	"net/http"
	"net/url"

	"github.com/fabmob/playground-standard-covoiturage/cmd/test"
	"github.com/spf13/cobra"
)

// driverJourneysCmd represents the driverJourneys command
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
		&departureDate, "departureTimeOfDay", "", "departureTimeOfDay query parameter")
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

func getDriverRegularTripsRun(runner test.TestRunner, server string) error {
	query := makeJourneyQuery()
	URL, _ := url.JoinPath(server, "/driver_regular_trips")

	return runner.Run(http.MethodGet, URL, verbose, query, nil, flags(http.StatusOK))
}

func makeRegularTripQuery() test.Query {
	var query = test.NewQuery()
	query.Params["departureLat"] = departureLat
	query.Params["departureLng"] = departureLng
	query.Params["arrivalLat"] = arrivalLat
	query.Params["arrivalLng"] = arrivalLng
	query.Params["departureDate"] = departureDate
	if timeDelta != "" {
		query.Params["timeDelta"] = timeDelta
	}
	if departureRadius != "" {
		query.Params["departureRadius"] = departureRadius
	}
	if arrivalRadius != "" {
		query.Params["arrivalRadius"] = arrivalRadius
	}
	if count != "" {
		query.Params["count"] = count
	}
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
