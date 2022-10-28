package test

import "fmt"

// Report stores assertionResults.
type Report struct {
	verbose          bool
	endpoint         Endpoint
	assertionResults []AssertionResult
}

func NewReport(assertionResults ...AssertionResult) Report {
	return Report{assertionResults: assertionResults}
}

func (report *Report) String() string {
	str := ""
	for _, ar := range report.assertionResults {
		if ar.Unwrap() == nil && report.verbose {
			str += stringOK(format(report.endpoint, ar.assertionDescription))
		} else if err := ar.Unwrap(); err != nil {
			str += stringError(format(report.endpoint, ar.assertionDescription))
			str += stringDetail(err.Error())
		}
	}
	return str
}

func (report *Report) countErrors() int {
	nErr := 0
	for _, ar := range report.assertionResults {
		if ar.Unwrap() != nil {
			nErr++
		}
	}
	return nErr
}

func format(endpoint Endpoint, assertionDescription string) string {
	return fmt.Sprintf("%-35s %-35s", endpoint, assertionDescription)

}

func (report *Report) hasErrors() bool {
	return report.countErrors() > 0
}
