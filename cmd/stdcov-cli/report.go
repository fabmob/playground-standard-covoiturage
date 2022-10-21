package main

import "fmt"

// Report stores assertionResults.
type Report struct {
	verbose             bool
	allAssertionResults []AssertionResult
}

func (report *Report) String() string {
	str := ""
	for _, ar := range report.allAssertionResults {
		str += toString(ar, report.verbose)
	}
	return str
}

func (report *Report) countErrors() int {
	nErr := 0
	for _, ar := range report.allAssertionResults {
		if ar.Unwrap() != nil {
			nErr++
		}
	}
	return nErr
}

func (report *Report) hasErrors() bool {
	return report.countErrors() > 0
}

func toString(ar AssertionResult, verbose bool) string {
	if shouldReport(ar, verbose) {
		return fmt.Sprintf("%s\n", ar.String())
	}
	return ""

}

func shouldReport(ar AssertionResult, verbose bool) bool {
	if ar.Unwrap() != nil || verbose {
		return true
	}
	return false
}
