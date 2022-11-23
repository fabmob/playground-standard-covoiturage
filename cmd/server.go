package cmd

import (
	"github.com/fabmob/playground-standard-covoiturage/cmd/service"
	"github.com/spf13/cobra"
)

// serveCmd represents the server command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serves a test API enforcing the standard covoitrage specification",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		service.Run(dataFile)
	},
}

var dataFile string

func init() {
	serveCmd.Flags().StringVar(&dataFile, "data", "", "Path to custom initial data file")

	rootCmd.AddCommand(serveCmd)
}
