package cmd

import (
	"net/http"
	"net/url"

	"github.com/fabmob/playground-standard-covoiturage/cmd/test"
	"github.com/spf13/cobra"
)

// passengerJourneysCmd represents the passengerJourneys command
var passengerJourneysCmd = makeEndpointCommand(test.GetPassengerJourneysEndpoint)

func init() {
	passengerJourneysCmd.PreRunE = checkGetJourneysCmdFlags
	passengerJourneysCmd.Run = func(cmd *cobra.Command, args []string) {
		query := makeJourneyQuery()
		URL, _ := url.JoinPath(server, "/passenger_journeys")
		err := test.RunTest(http.MethodGet, URL, verbose, query, nil, flags(http.StatusOK))
		exitWithError(err)
	}

	passengerJourneysCmd.Flags().StringVar(&departureLat, "departureLat", "", "DepartureLat parameters in the form name = value")
	passengerJourneysCmd.Flags().StringVar(&departureLng, "departureLng", "", "DepartureLng parameters in the form name = value")
	passengerJourneysCmd.Flags().StringVar(&arrivalLat, "arrivalLat", "", "ArrivalLat parameters in the form name = value")
	passengerJourneysCmd.Flags().StringVar(&arrivalLng, "arrivalLng", "", "ArrivalLng parameters in the form name = value")
	passengerJourneysCmd.Flags().StringVar(&departureDate, "departureDate", "", "DepartureDate parameters in the form name = value")
	passengerJourneysCmd.Flags().StringVar(&timeDelta, "timeDelta", "", "TimeDelta parameters in the form name = value")
	passengerJourneysCmd.Flags().StringVar(&departureRadius, "departureRadius", "", "DepartureRadius parameters in the form name = value")
	passengerJourneysCmd.Flags().StringVar(&arrivalRadius, "arrivalRadius", "", "ArrivalRadius parameters in the form name = value")
	passengerJourneysCmd.Flags().StringVar(&count, "count", "", "Count parameters in the form name = value")

	getCmd.AddCommand(passengerJourneysCmd)
}