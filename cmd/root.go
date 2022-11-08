package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pscovoit",
	Short: "Testing tool for the standard-covoiturage API",
	Long: `This testing tool can be used to serve a test API enforcing the
  standard covoiturage specification, for exploration and experimentation, or to run a test
  suite against a request to a custom API.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	exitWithError(err)
}

func init() {}

func exitWithError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
