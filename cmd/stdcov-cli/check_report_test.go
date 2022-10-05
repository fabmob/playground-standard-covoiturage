package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCheckReport_String(t *testing.T) {
	testCases := []struct {
		statusCheck    bool
		expectedString string
	}{
		{false, "❌ Check of the status endpoint failed\n"},
		{true, "✅ Status check OK\n"},
	}
	for _, tc := range testCases {
		cr := CheckReport{statusCheck: tc.statusCheck}
		crStr := cr.String()
		if crStr != tc.expectedString {
			t.Errorf(
				"Report to string does not produce the expected message.\nDiff: %s",
				cmp.Diff(crStr, tc.expectedString),
			)
		}
	}
}
