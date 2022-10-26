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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
