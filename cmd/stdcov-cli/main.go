package main

import (
	"flag"
	"fmt"
	"os"

	"gitlab.com/multi/stdcov-api-test/cmd/stdcov-cli/client"
)

//go:generate oapi-codegen -package client -o ./client/client.go -generate "types,client" --old-config-style stdcov_openapi.yaml

func main() {
	os.Exit(run())
}

// run runs the validation and returns an exit code
func run() int {
	urlStrPtr := flag.String("url", "", "Base url of the API under test")
	verboseBoolPtr := flag.Bool("verbose", false, "Make the operation more talkative")

	flag.Parse()

	client, err := client.NewClient(*urlStrPtr)
	if err != nil {
		panic(err)
	}
	report := ExecuteTestSuite(client)
	report.verbose = *verboseBoolPtr
	fmt.Println(report)
	if report.hasErrors() {
		return 1
	}
	return 0
}
