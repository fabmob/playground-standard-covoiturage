package test

import (
	"fmt"
	"net/http"
	"net/url"

	"gitlab.com/multi/stdcov-api-test/cmd/api"
)

// Run runs the cli validation and returns an exit code
func Run(server, URL string, verbose bool, query Query) int {

	c, err := api.NewClient(server)
	if err != nil {
		panic(err)
	}
	fullURL, _ := url.JoinPath(server, URL)
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		fmt.Println(err)
		return 1
	}

	AddQueryParameters(query, req)

	flags := Flags{
		DisallowEmpty: false,
	}
	report, err := Request(c, req, flags)
	if err != nil {
		fmt.Println(err)
		return 1
	}

	report.verbose = verbose
	fmt.Println(report)
	if report.hasErrors() {
		return 1
	}
	return 0
}
