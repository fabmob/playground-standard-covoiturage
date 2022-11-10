package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:               "get",
	Short:             "Interface for testing endpoints with method GET",
	Long:              "Interface for testing endpoints with method GET",
	PersistentPreRunE: checkCmdFlags,
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

func initMethodFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&server, "server", "", "Server on which torun the query")
}

func initGetCmd() {
	initMethodFlags(getCmd)
	testCmd.AddCommand(getCmd)
}

func checkCmdFlags(cmd *cobra.Command, args []string) error {
	if server == "" {
		return errors.New("missing required --server information")
	}

	return nil
}

// postCmd represents the post command
var postCmd = &cobra.Command{
	Use:               "post",
	Short:             "Interface for testing endpoints with method POST",
	Long:              "Interface for testing endpoints with method POST",
	PersistentPreRunE: checkCmdFlags,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("post called")
	},
}

func initPostCmd() {
	initMethodFlags(postCmd)
	testCmd.AddCommand(postCmd)
}

func init() {
	initGetCmd()
	initPostCmd()
}
