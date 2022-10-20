package main

import (
	"errors"
	"net/http"
	"net/url"
	"testing"
)

// throwErrorOnTest returns an urlError for every API call and checks:
// - that only one AssertionError is returned
// - that AssertionError.Unwrap() != nil
func testThrowErrorOnTest(t *testing.T, f auxTestFun, r *http.Request) {
	t.Run("API call throws error", func(t *testing.T) {
		urlError := &url.Error{Op: "", URL: "", Err: errors.New("error")}
		m := returnErrorClient(urlError)
		a := NewAssertionAccu()

		f(m, r, a)
		shouldHaveSingleAssertionResult(t, a)

		err := a.GetAssertionResults()[0].Unwrap()
		if err == nil {
			t.Error("If error returned, api is not up")
		}
	})
}

func TestAPIErrors(t *testing.T) {
	testFuns := []struct {
		f auxTestFun
		r *http.Request
	}{
		testGetStatus,
		testGetDriverJourneys,
	}
	for _, f := range testFuns {
		testThrowErrorOnTest(t, f)
	}
}
