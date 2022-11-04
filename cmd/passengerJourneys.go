package cmd

import (
	"fmt"
	"time"

	"github.com/fabmob/playground-standard-covoiturage/cmd/test"
	"github.com/spf13/cobra"
)

// passengerJourneysCmd represents the passengerJourneys command
var passengerJourneysCmd = &cobra.Command{
	Use:   "passengerJourneys",
	Short: "Test the GET /passenger_journeys endpoint",
	Long: `Test the GET /passenger_journeys endpoint

Default query coordinates are placed on "Vesdun", a small town proclaimed "center
of France" by IGN in 1993.`,
	Run: func(cmd *cobra.Command, args []string) {
		query := makeQuery()
		test.Run(server, "/passenger_journeys", verbose, query)
	},
}

func init() {
	driverJourneysCmd.Flags().StringVar(&departureLat, "departureLat", vesdunLat, "DepartureLat parameters in the form name = value")
	driverJourneysCmd.Flags().StringVar(&departureLng, "departureLng", vesdunLng, "DepartureLng parameters in the form name = value")
	driverJourneysCmd.Flags().StringVar(&arrivalLat, "arrivalLat", vesdunLat, "ArrivalLat parameters in the form name = value")
	driverJourneysCmd.Flags().StringVar(&arrivalLng, "arrivalLng", vesdunLng, "ArrivalLng parameters in the form name = value")
	driverJourneysCmd.Flags().StringVar(&departureDate, "departureDate", fmt.Sprintf("%d", time.Now().Unix()), "DepartureDate parameters in the form name = value")
	driverJourneysCmd.Flags().StringVar(&timeDelta, "timeDelta", "", "TimeDelta parameters in the form name = value")
	driverJourneysCmd.Flags().StringVar(&departureRadius, "departureRadius", "", "DepartureRadius parameters in the form name = value")
	driverJourneysCmd.Flags().StringVar(&arrivalRadius, "arrivalRadius", "", "ArrivalRadius parameters in the form name = value")
	driverJourneysCmd.Flags().StringVar(&count, "count", "", "Count parameters in the form name = value")

	getCmd.AddCommand(passengerJourneysCmd)
}
