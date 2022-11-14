package test

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type TestRunner interface {
	Run(method, URL string, verbose bool, query Query, body []byte, flags Flags) error
}

type DefaultRunner struct{}

// Run runs the cli validation and returns an exit code
func (DefaultRunner) Run(method, URL string, verbose bool, query Query, body []byte, flags Flags) error {

	req, err := http.NewRequest(method, URL, nil)
	if err != nil {
		return err
	}

	AddQueryParameters(query, req)

	report, err := Request(req, flags)
	if err != nil {
		return err
	}

	report.verbose = verbose
	fmt.Println(report)

	if report.hasErrors() {
		return errors.New(report.String())
	}

	return nil
}

func RunTest(method, URL string, verbose bool, query Query, body []byte, flags Flags) error {
	runner := DefaultRunner{}
	return runner.Run(method, URL, verbose, query, body, flags)
}
