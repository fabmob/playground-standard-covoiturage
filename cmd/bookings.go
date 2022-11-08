package cmd

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/fabmob/playground-standard-covoiturage/cmd/test"
	"github.com/spf13/cobra"
)

// bookingsCmd represents the bookings command
var bookingsCmd = &cobra.Command{
	Use:   "bookings",
	Short: "Test the GET /bookings endpoint",
	Long:  `Test the GET /bookings endpoint`,
	Run: func(cmd *cobra.Command, args []string) {
		checkRequiredGetBookingsArguments(server, bookingID)
		URL, _ := url.JoinPath(server, "/bookings", bookingID)
		test.Run(http.MethodGet, URL, verbose, test.NewQuery(), flags())
	},
}

var (
	bookingID string
)

func init() {
	bookingsCmd.Flags().StringVar(
		&bookingID, "bookingId", "", "bookingIdpath parameter",
	)
	getCmd.AddCommand(bookingsCmd)
}

func checkRequiredGetBookingsArguments(server, bookingID string) {
	if server == "" {
		fmt.Println("Server information is missing, but required. It can be passed with the --server flag.")
		os.Exit(0)
	}
	if bookingID == "" {
		fmt.Println("bookingId path parameter is missing, but required. It can be passed with the --bookingId flag.")
		os.Exit(0)
	}
}
