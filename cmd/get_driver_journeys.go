package cmd

import (
	"net/http"
	"net/url"

	"github.com/fabmob/playground-standard-covoiturage/cmd/test"
	"github.com/spf13/cobra"
)

// driverJourneysCmd represents the driverJourneys command
var driverJourneysCmd = makeEndpointCommand(test.GetDriverJourneysEndpoint)

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

func init() {

	driverJourneysCmd.PreRunE = checkGetJourneysCmdFlags

	driverJourneysCmd.Run = func(cmd *cobra.Command, args []string) {
		query := makeJourneyQuery(departureLat, departureLng, arrivalLat, arrivalLng, departureDate, timeDelta, departureRadius, arrivalRadius, count)
		URL, _ := url.JoinPath(server, "/driver_journeys")
		err := test.RunTest(http.MethodGet, URL, verbose, query, nil, flags(http.StatusOK))
		exitWithError(err)
	}

	driverJourneysCmd.Flags().StringVar(
		&departureLat, "departureLat", "", "departureLat query query parameter")
	driverJourneysCmd.Flags().StringVar(
		&departureLng, "departureLng", "", "departureLng query parameter")
	driverJourneysCmd.Flags().StringVar(
		&arrivalLat, "arrivalLat", "", "arrivalLat query parameter")
	driverJourneysCmd.Flags().StringVar(
		&arrivalLng, "arrivalLng", "", "arrivalLng query parameter")
	driverJourneysCmd.Flags().StringVar(
		&departureDate, "departureDate", "", "departureDate query parameter")
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

func makeJourneyQuery(departureLat, departureLng, arrivalLat, arrivalLng, departureDate, timeDelta, departureRadius, arrivalRadius, count string) test.Query {
	var query = test.NewQuery()

	query.SetParam("departureLat", departureLat)
	query.SetParam("departureLng", departureLng)
	query.SetParam("arrivalLat", arrivalLat)
	query.SetParam("arrivalLng", arrivalLng)
	query.SetParam("departureDate", departureDate)

	query.SetOptionalParam("timeDelta", timeDelta)
	query.SetOptionalParam("departureRadius", departureRadius)
	query.SetOptionalParam("arrivalRadius", arrivalRadius)
	query.SetOptionalParam("count", count)

	return query
}

func checkGetJourneysCmdFlags(cmd *cobra.Command, args []string) error {
	return anyError(
		checkRequiredDepartureLat(departureLat),
		checkRequiredDepartureLng(departureLng),
		checkRequiredArrivalLat(departureLat),
		checkRequiredArrivalLng(departureLng),
		checkRequiredDepartureDate(departureDate),
	)
}
