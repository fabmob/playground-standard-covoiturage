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
		if backgroundProcess {
			go service.Run()
		} else {
			service.Run()
		}
	},
}

var backgroundProcess bool

func init() {
	serveCmd.Flags().BoolVar(&backgroundProcess, "bg", false, "Run the server in a background process")
	rootCmd.AddCommand(serveCmd)
}
