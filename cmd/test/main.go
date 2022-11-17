package test

import (
	"fmt"
)

type TestRunner interface {
	Run(method, URL string, verbose bool, query Query, body []byte, apiKey string, flags Flags) error
}

type DefaultRunner struct{}

func NewDefaultRunner() *DefaultRunner {
	return &DefaultRunner{}
}

// Run runs the cli validation and returns an exit code
func (*DefaultRunner) Run(method, URL string, verbose bool, query Query, body []byte, apiKey string, flags Flags) error {

	req, err := makeRequest(method, URL, body, apiKey)
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
		return fmt.Errorf("❌ %d failed assertion(s) ", report.countErrors())
	}

	return nil
}

func RunTest(method, URL string, verbose bool, query Query, body []byte, apiKey string, flags Flags) error {
	runner := DefaultRunner{}
	return runner.Run(method, URL, verbose, query, body, apiKey, flags)
}
