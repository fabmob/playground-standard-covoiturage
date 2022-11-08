package test

import (
	"fmt"
	"net/http"
)

// Run runs the cli validation and returns an exit code
func Run(method, URL string, verbose bool, query Query, flags Flags) int {

	initAPIMapping()

	req, err := http.NewRequest(method, URL, nil)
	if err != nil {
		fmt.Println(err)
		return 1
	}

	AddQueryParameters(query, req)

	report, err := Request(req, flags)
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
