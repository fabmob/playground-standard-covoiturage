package cmd

import (
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
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.stdcov-api-test.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	/* rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle") */
}
