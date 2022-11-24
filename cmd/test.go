package cmd

import (
	"net/http"
	"time"

	"github.com/fabmob/playground-standard-covoiturage/cmd/test"
	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Test an API complying with the standard covoiturage",
	Long:  "Test an API complying with the standard covoiturage",
	Run: func(cmd *cobra.Command, args []string) {
		var timeout = 100 * time.Millisecond

		body, err := readBodyFromStdin(cmd, timeout)
		if err != nil {
			body = nil
		}

		err = test.RunTest(method, URL, verbose, query, body, apiKey, flags(http.StatusOK))
		exitWithError(err)
	},
}

var (
	apiKey        string
	URL           string
	verbose       bool
	query         test.Query
	disallowEmpty bool
	expectStatus  int
	method        string
)

func init() {
	rootCmd.AddCommand(testCmd)

	testCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Make the operation more talkative")
	testCmd.PersistentFlags().BoolVar(
		&disallowEmpty,
		"disallowEmpty",
		test.DefaultDisallowEmptyFlag,
		"Should an empty request return an error",
	)
	testCmd.PersistentFlags().StringVar(&apiKey, "auth", "", "API key sent in the \"X-API-Key\" header of the request")
	testCmd.PersistentFlags().IntVar(
		&expectStatus,
		"expectStatus",
		0,
		"Expected status code. Defaults to success, 2xx, status code - exact default depends on endpoint",
	)

	testCmd.Flags().StringVar(&method, "method", http.MethodGet, "HTTP method, either GET (default), POST or PATCH")
	testCmd.Flags().StringVarP(&URL, "url", "u", "", "API call URL")
	testCmd.Flags().VarP(&query, "query", "q", "Query parameters in the form name=value")
}

func flags(defaultStatus int) test.Flags {
	flags := test.NewFlags()
	flags.DisallowEmpty = disallowEmpty
	if expectStatus == 0 { //not set
		flags.ExpectedStatusCode = defaultStatus
	} else {
		flags.ExpectedStatusCode = expectStatus
	}
	return flags
}
