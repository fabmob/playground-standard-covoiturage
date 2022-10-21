package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"gitlab.com/multi/stdcov-api-test/cmd/stdcov-cli/client"
)

//go:generate oapi-codegen -package client -o ./client/client.go -generate "types,client" --old-config-style ../../spec/stdcov_openapi.yaml

func main() {
	os.Exit(Run())
}

// Run runs the cli validation and returns an exit code
func Run() int {
	serverStrPtr := flag.String("server", "", "Server URL of the API under test")
	urlStrPtr := flag.String("url", "", "API call URL")
	verboseBoolPtr := flag.Bool("verbose", false, "Make the operation more talkative")
	var query Query
	flag.Var(&query, "q", "Query parameters in the form name=value")

	flag.Parse()

	c, err := client.NewClient(*serverStrPtr)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("GET", *urlStrPtr, nil)
	if err != nil {
		fmt.Println(err)
		return 1
	}

	AddQueryParameters(query, req)

	report, err := ExecuteTestSuite(c, req)
	if err != nil {
		fmt.Println(err)
		return 1
	}

	report.verbose = *verboseBoolPtr
	fmt.Println(report)
	if report.hasErrors() {
		return 1
	}
	return 0
}
