package cmd

import (
	"github.com/fabmob/playground-standard-covoiturage/cmd/endpoint"
	"github.com/fabmob/playground-standard-covoiturage/cmd/test"
	"github.com/spf13/cobra"
)

var passengerRegularTripsCmd = makeEndpointCommand(endpoint.GetPassengerRegularTrips)

var getPassengerRegularTripsParameters = getDriverRegularTripsParameters

func init() {
	cmd := passengerRegularTripsCmd
	cmd.PreRunE = checkGetRegularTripsCmdFlags

	cmd.Run = func(cmd *cobra.Command, args []string) {
		err := getRegularTripsRun(
			test.NewDefaultRunner(),
			server,
			getPassengerRegularTripsParameters,
			departureWeekdays,
			"/passenger_regular_trips",
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
