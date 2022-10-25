package stdcovcli

import (
	"errors"
	"testing"
)

func TestReportSingle(t *testing.T) {
	testCases := []struct {
		name         string
		err          error
		verbose      bool
		shouldReport bool
	}{
		{
			"No error, verbose false",
			nil,
			false,
			false,
		},
		{
			"No error, verbose true",
			nil,
			true,
			true,
		},
		{
			"Error, verbose false",
			errors.New(""),
			false,
			true,
		},
		{
			"Error, verbose true",
			errors.New(""),
			true,
			true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ar := NewAssertionResult(tc.err, "", "", "")
			if shouldReport(ar, tc.verbose) != tc.shouldReport {
				t.Logf("verbose: %t", tc.verbose)
				t.Logf("hasError: %t", tc.err == nil)
				t.Error("Report single has wrong behavior")
			}
		})
	}
}

func TestReportCountErrors(t *testing.T) {
	testCases := []struct {
		allAssertionErr []error
		expectedNErr    int
	}{
		{[]error{}, 0},
		{[]error{nil}, 0},
		{[]error{errors.New("")}, 1},
		{[]error{errors.New(""), nil}, 1},
		{[]error{errors.New(""), errors.New("")}, 2},
		{[]error{nil, errors.New(""), errors.New(""), nil}, 2},
	}
	for _, tc := range testCases {
		assertionResults := []AssertionResult{}
		for _, err := range tc.allAssertionErr {
			assertionResults = append(assertionResults, NewAssertionResult(err, "", "", ""))
		}
		report := Report{false, assertionResults}

		if report.countErrors() != tc.expectedNErr {
			t.Error("Wrong number of errors")
		}

		expectedHasErr := tc.expectedNErr > 0
		if report.hasErrors() != expectedHasErr {
			t.Error("Report.hasError does not correctly report if report has errors")
		}
	}
}
