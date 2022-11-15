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
		query := makeJourneyQuery()
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

func makeJourneyQuery() test.Query {
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

func checkGetJourneysCmdFlags(cmd *cobra.Command, args []string) error {
	return anyError(
		checkRequiredDepartureLat(departureLat),
		checkRequiredDepartureLng(departureLng),
		checkRequiredArrivalLat(departureLat),
		checkRequiredArrivalLng(departureLng),
		checkRequiredDepartureDate(departureDate),
	)
}
