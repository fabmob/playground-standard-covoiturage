package test

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/fabmob/playground-standard-covoiturage/cmd/endpoint"
	"github.com/fabmob/playground-standard-covoiturage/cmd/test/assert"
)

func TestReport(t *testing.T) {
	endpoint := endpoint.New(http.MethodGet, "/endpoint_path")
	assertStr := "test assertion"
	errorDescription := "Error description"

	makeReport := func(err error, verbose bool) Report {
		report := NewReport(assert.NewAssertionResult(err, assertStr))
		report.endpoint = endpoint
		report.verbose = verbose
		return report
	}

	shouldContain := func(t *testing.T, r Report, str string) {
		t.Helper()
		if !strings.Contains(r.String(), str) {
			t.Logf("Assertion string : %s", r.String())
			t.Error("Assertion string does not contain " + str)
		}
	}

	testCases := []struct {
		name    string
		err     error
		verbose bool
	}{
		{
			"Assertion without error",
			nil,
			true,
		},
		{
			"Assertion with error",
			errors.New(errorDescription),
			true,
		},
	}

	for _, tc := range testCases {
		t.Run("Assertion with error", func(t *testing.T) {
			r := makeReport(tc.err, tc.verbose)
			shouldContain(t, r, endpoint.Method)
			shouldContain(t, r, endpoint.Path)
			shouldContain(t, r, assertStr)

			if tc.err != nil {
				shouldContain(t, r, errorDescription)
			}
		})
	}
}

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
			ar := assert.NewAssertionResult(tc.err, "")
			report := NewReport(ar)
			report.verbose = tc.verbose

			if (report.String() != "") != tc.shouldReport {
				t.Logf("verbose: %t", tc.verbose)
				t.Logf("hasError: %t", tc.err == nil)
				t.Logf("report: %s", report.String())
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
		var assertionResults = []assert.Result{}

		for _, err := range tc.allAssertionErr {
			assertionResults = append(assertionResults, assert.NewAssertionResult(err, ""))
		}

		report := NewReport(assertionResults...)

		if report.countErrors() != tc.expectedNErr {
			t.Error("Wrong number of errors")
		}

		expectedHasErr := tc.expectedNErr > 0
		if report.hasErrors() != expectedHasErr {
			t.Error("Report.hasError does not correctly report if report has errors")
		}
	}
}
