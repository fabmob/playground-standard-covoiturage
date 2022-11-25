package cmd

import (
	"net/http"
	"net/url"

	"github.com/fabmob/playground-standard-covoiturage/cmd/endpoint"
	"github.com/fabmob/playground-standard-covoiturage/cmd/test"
	"github.com/spf13/cobra"
)

var passengerRegularTripsCmd = makeEndpointCommand(endpoint.GetPassengerRegularTrips)

func init() {
	passengerRegularTripsCmd.PreRunE = checkGetRegularTripsCmdFlags

	passengerRegularTripsCmd.Run = func(cmd *cobra.Command, args []string) {
		err := getPassengerRegularTripsRun(
			test.NewDefaultRunner(),
			server,
			departureLat, departureLng, arrivalLat, arrivalLng, departureTimeOfDay,
			departureWeekdays, timeDelta, departureRadius, arrivalRadius, count,
			minDepartureDate, maxDepartureDate,
		)
		exitWithError(err)
	}

	passengerRegularTripsCmd.Flags().StringVar(
		&departureLat, "departureLat", "", "departureLat query query parameter")
	passengerRegularTripsCmd.Flags().StringVar(
		&departureLng, "departureLng", "", "departureLng query parameter")
	passengerRegularTripsCmd.Flags().StringVar(
		&arrivalLat, "arrivalLat", "", "arrivalLat query parameter")
	passengerRegularTripsCmd.Flags().StringVar(
		&arrivalLng, "arrivalLng", "", "arrivalLng query parameter")
	passengerRegularTripsCmd.Flags().StringVar(
		&departureDate, "departureTimeOfDay", "", "departureTimeOfDay query parameter")
	passengerRegularTripsCmd.Flags().StringSliceVar(
		&departureWeekdays, "departureWeekdays", []string{}, "departureWeekdays query parameter")
	passengerRegularTripsCmd.Flags().StringVar(
		&timeDelta, "timeDelta", "", "timeDelta query parameter")
	passengerRegularTripsCmd.Flags().StringVar(
		&departureRadius, "departureRadius", "", "departureRadius query parameter")
	passengerRegularTripsCmd.Flags().StringVar(
		&arrivalRadius, "arrivalRadius", "", "arrivalRadius query parameter")
	passengerRegularTripsCmd.Flags().StringVar(
		&minDepartureDate, "minDepartureDate", "", "minDepartureDate query parameter")
	passengerRegularTripsCmd.Flags().StringVar(
		&maxDepartureDate, "maxDepartureDate", "", "maxDepartureDate query parameter")
	passengerRegularTripsCmd.Flags().StringVar(
		&count, "count", "", "count query parameter")

	getCmd.AddCommand(passengerRegularTripsCmd)
}

func getPassengerRegularTripsRun(
	runner test.TestRunner,
	server,
	departureLat, departureLng, arrivalLat, arrivalLng, departureTimeOfDay string,
	departureWeekdays []string,
	timeDelta, departureRadius, arrivalRadius, count, minDepartureDate,
	maxDepartureDate string,
) error {
	query := makeRegularTripQuery(
		departureLat, departureLng, arrivalLat, arrivalLng, departureTimeOfDay,
		departureWeekdays, timeDelta, departureRadius, arrivalRadius, count,
		minDepartureDate, maxDepartureDate,
	)
	URL, _ := url.JoinPath(server, "/passenger_regular_trips")

	return runner.Run(http.MethodGet, URL, query, nil, verbose, apiKey, flagsWithDefault(http.StatusOK))
}
