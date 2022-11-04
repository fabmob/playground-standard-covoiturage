package cmd

import (
	"fmt"
	"net/http"
	"net/url"
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
		URL, _ := url.JoinPath(server, "/passenger_journeys")
		test.Run(http.MethodGet, URL, verbose, query, flags())
	},
}

func init() {
	passengerJourneysCmd.Flags().StringVar(&departureLat, "departureLat", vesdunLat, "DepartureLat parameters in the form name = value")
	passengerJourneysCmd.Flags().StringVar(&departureLng, "departureLng", vesdunLng, "DepartureLng parameters in the form name = value")
	passengerJourneysCmd.Flags().StringVar(&arrivalLat, "arrivalLat", vesdunLat, "ArrivalLat parameters in the form name = value")
	passengerJourneysCmd.Flags().StringVar(&arrivalLng, "arrivalLng", vesdunLng, "ArrivalLng parameters in the form name = value")
	passengerJourneysCmd.Flags().StringVar(&departureDate, "departureDate", fmt.Sprintf("%d", time.Now().Unix()), "DepartureDate parameters in the form name = value")
	passengerJourneysCmd.Flags().StringVar(&timeDelta, "timeDelta", "", "TimeDelta parameters in the form name = value")
	passengerJourneysCmd.Flags().StringVar(&departureRadius, "departureRadius", "", "DepartureRadius parameters in the form name = value")
	passengerJourneysCmd.Flags().StringVar(&arrivalRadius, "arrivalRadius", "", "ArrivalRadius parameters in the form name = value")
	passengerJourneysCmd.Flags().StringVar(&count, "count", "", "Count parameters in the form name = value")

	getCmd.AddCommand(passengerJourneysCmd)
}
