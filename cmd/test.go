package cmd

import (
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
		test.Run(server, url, verbose, query)
	},
}

var (
	server        string
	url           string
	verbose       bool
	query         test.Query
	disallowEmpty bool
)

func init() {
	rootCmd.AddCommand(testCmd)
	testCmd.PersistentFlags().StringVarP(&server, "server", "s", "", "Server URL of the API under test")
	testCmd.PersistentFlags().StringVarP(&url, "url", "u", "", "API call URL")
	testCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Make the operation more talkative")
	testCmd.PersistentFlags().VarP(&query, "query", "q", "Query parameters in the form name=value")
	testCmd.PersistentFlags().BoolVar(&disallowEmpty, "disallowEmpty", false,
		"Should an empty request return an error")
}
