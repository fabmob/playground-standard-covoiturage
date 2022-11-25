package cmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	server string
)

var (
	getCmd   = methodCmdHelper("get")
	postCmd  = methodCmdHelper("post")
	patchCmd = methodCmdHelper("patch")
)

func init() {
	initMethodCmd(getCmd)
	initMethodCmd(postCmd)
	initMethodCmd(patchCmd)
}

func methodCmdHelper(method string) *cobra.Command {
	description := "Interface for testing endpoints with method " + strings.ToUpper(method)
	return &cobra.Command{
		Use:               method,
		Short:             description,
		Long:              description,
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
}

func initMethodCmd(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&server, "server", "", "(required) Server on which torun the query")
	testCmd.AddCommand(cmd)
}

func checkCmdFlags(cmd *cobra.Command, args []string) error {
	return checkRequiredServer(server)
}
