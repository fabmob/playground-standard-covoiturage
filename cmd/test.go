package cmd

import (
	"net/http"

	"github.com/fabmob/playground-standard-covoiturage/cmd/test"
	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := test.Run(http.MethodGet, URL, verbose, query, flags())
		exitWithError(err)
	},
}

var (
	URL           string
	verbose       bool
	query         test.Query
	disallowEmpty bool
	expectStatus  int
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
	testCmd.PersistentFlags().IntVar(
		&expectStatus,
		"expectStatus",
		test.DefaultExpectedStatusCode,
		"Expected status code",
	)

	testCmd.Flags().StringVarP(&URL, "url", "u", "", "API call URL")
	testCmd.Flags().VarP(&query, "query", "q", "Query parameters in the form name=value")
}

func flags() test.Flags {
	flags := test.NewFlags()
	flags.DisallowEmpty = disallowEmpty
	flags.ExpectedStatusCode = expectStatus
	return flags
}
