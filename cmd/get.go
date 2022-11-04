package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Interface for testing endpoints with method GET",
	Long:  "",
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

func init() {
	testCmd.AddCommand(getCmd)
}
