package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/multi/stdcov-api-test/cmd/service"
)

// serveCmd represents the server command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serves a test API enforcing the standard covoitrage specification",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		service.Run()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
