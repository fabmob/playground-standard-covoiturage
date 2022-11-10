package cmd

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/fabmob/playground-standard-covoiturage/cmd/test"
	"github.com/spf13/cobra"
)

// driverJourneysCmd represents the driverJourneys command
var driverJourneysCmd = &cobra.Command{
	Use:   "driverJourneys",
	Short: "Test the GET /driver_journeys endpoint",
	Long: `Test the GET /driver_journeys endpoint

Default query coordinates are placed on "Vesdun", a small town proclaimed "center
of France" by IGN in 1993.`,
	Run: func(cmd *cobra.Command, args []string) {
		query := makeQuery()
		URL, _ := url.JoinPath(server, "/driver_journeys")
		err := test.Run(http.MethodGet, URL, verbose, query, nil, flags(http.StatusOK))
		exitWithError(err)
	},
}

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

const (
	vesdunLat = "46.588"
	vesdunLng = "2.4284"
)

func init() {
	driverJourneysCmd.Flags().StringVar(
		&departureLat, "departureLat", vesdunLat, "departureLat query query parameter")
	driverJourneysCmd.Flags().StringVar(
		&departureLng, "departureLng", vesdunLng, "departureLng query parameter")
	driverJourneysCmd.Flags().StringVar(
		&arrivalLat, "arrivalLat", vesdunLat, "arrivalLat query parameter")
	driverJourneysCmd.Flags().StringVar(
		&arrivalLng, "arrivalLng", vesdunLng, "arrivalLng query parameter")
	driverJourneysCmd.Flags().StringVar(
		&departureDate, "departureDate", fmt.Sprintf("%d", time.Now().Unix()), "departureDate query parameter")
	driverJourneysCmd.Flags().StringVar(
		&timeDelta, "timeDelta", "", "timeDelta query parameter")
	driverJourneysCmd.Flags().StringVar(
		&departureRadius, "departureRadius", "", "departureRadius query parameter")
	driverJourneysCmd.Flags().StringVar(
		&arrivalRadius, "arrivalRadius", "", "arrivalRadius query parameter")
	driverJourneysCmd.Flags().StringVar(
		&count, "count", "", "count query parameter")

	getCmd.AddCommand(driverJourneysCmd)
}

func makeQuery() test.Query {
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
