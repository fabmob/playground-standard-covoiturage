package test

import (
	"fmt"

	"github.com/fabmob/playground-standard-covoiturage/cmd/endpoint"
	"github.com/fabmob/playground-standard-covoiturage/cmd/test/assert"
)

// Report stores and prints `assert.Result`s
type Report struct {
	verbose          bool
	endpoint         endpoint.Info
	assertionResults []assert.Result
}

// NewReport creates a new report with given assertion results
func NewReport(assertionResults ...assert.Result) Report {
	return Report{assertionResults: assertionResults}
}

// String implements stringer interface. It is
// used to print the report to terminal
func (report *Report) String() string {
	str := ""

	for _, ar := range report.assertionResults {
		if ar.Unwrap() == nil && report.verbose {
			str += stringOK(format(report.endpoint, ar.AssertionDescription))
		} else if err := ar.Unwrap(); err != nil {
			str += stringError(format(report.endpoint, ar.AssertionDescription))
			str += stringDetail(err.Error())
		}
	}

	return str
}

func (report *Report) countErrors() int {
	var nErr = 0

	for _, ar := range report.assertionResults {
		if ar.Unwrap() != nil {
			nErr++
		}
	}

	return nErr
}

func format(endpoint endpoint.Info, assertionDescription string) string {
	return fmt.Sprintf("%-35s %-35s", endpoint, assertionDescription)
}

func (report *Report) hasErrors() bool {
	return report.countErrors() > 0
}

// ////////////////////////////////////////////////////////////
// Printing helper functions
// ////////////////////////////////////////////////////////////

func stringError(msg string) string {
	return stringWithSymbol("ERROR ❌", msg)
}

func stringOK(msg string) string {
	return stringWithSymbol("OK ✅", msg)
}

func stringDetail(msg string) string {
	return stringWithSymbol("", msg)
}

func stringWithSymbol(symbol, msg string) string {
	return fmt.Sprintf(
		"%7s %s\n",
		symbol,
		msg,
	)
}
