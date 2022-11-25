// Package test has all utilities to perform a test on an endpoint of
// Standard-Covoiturage.
//
// Function `Run` runs adequate tests against a request to the standard
// covoiturage, detailed with method,
// URL, query and body.
//
// Test functions specific to a given endpoint are exported as
// `Test{Method}{Path}Response` functions (all camel case).
//
// The expectations of the tests can be refined thanks to the test `Flags`.
//
// `Assert*` functions are atomic assertions used in tests that can be extended through the
// `Assertion` interface.
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

	req, err := makeRequestWithContext(method, URL, body, apiKey)
	if err != nil {
		return err
	}

	AddQueryParameters(query, req)

	report, err := testRequest(req, flags)
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
