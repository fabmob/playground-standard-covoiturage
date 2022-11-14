package cmd

import (
	"net/http"
	"net/url"
	"time"

	"github.com/fabmob/playground-standard-covoiturage/cmd/test"
	"github.com/spf13/cobra"
)

var postBookingsCmd = &cobra.Command{
	Use:   "bookings",
	Short: "Test the POST /bookings endpoint",
	Long:  `Test the POST /bookings endpoint`,
	Run: func(cmd *cobra.Command, args []string) {

		var timeout = 100 * time.Millisecond

		body, err := readBodyFromStdin(cmd, timeout)
		exitWithError(err)

		URL, _ := url.JoinPath(server, "/bookings")
		err = test.RunTest(http.MethodPost, URL, verbose, test.NewQuery(), body, flags(http.StatusCreated))
		exitWithError(err)

	},
}

func init() {
	postCmd.AddCommand(postBookingsCmd)
}
