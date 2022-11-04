package test

import (
	"fmt"
	"net/http"
	"net/url"
)

// Run runs the cli validation and returns an exit code
func Run(server, endpoint string, verbose bool, query Query) int {

	fullURL, _ := url.JoinPath(server, endpoint)

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		fmt.Println(err)
		return 1
	}

	AddQueryParameters(query, req)

	flags := Flags{
		DisallowEmpty: false,
	}

	registerAllTests()

	report, err := Request(server, req, flags)
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
