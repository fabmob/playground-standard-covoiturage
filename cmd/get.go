package cmd

import (
	"errors"
	"os"

	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:               "get",
	Short:             "Interface for testing endpoints with method GET",
	Long:              "",
	PersistentPreRunE: checkGetCmdFlags,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			err := cmd.Help()
			if err != nil {
				panic(err)
			}

			os.Exit(0)
		}
	},
}

var (
	server string
)

func init() {
	getCmd.PersistentFlags().StringVar(&server, "server", "", "Server on which to run the query")
	testCmd.AddCommand(getCmd)
}

func checkGetCmdFlags(cmd *cobra.Command, args []string) error {
	if server == "" {
		return errors.New("missing required --server information")
	}

	return nil
}
