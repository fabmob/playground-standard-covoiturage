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
func testThrowErrorOnTest(t *testing.T, f auxTestFun) {
	t.Run("API call throws error", func(t *testing.T) {
		urlError := &url.Error{Op: "", URL: "", Err: errors.New("error")}
		m := returnErrorClient(urlError)
		a := NewAssertionAccu()

		// specific request is irrelevant as the error client is in any case returning
		// an error
		var r *http.Request
		f(m, r, a)
		shouldHaveSingleAssertionResult(t, a)

		err := a.GetAssertionResults()[0].Unwrap()
		if err == nil {
			t.Error("If error returned, api is not up")
		}
	})
}

func TestAPIErrors(t *testing.T) {

	testCases := []auxTestFun{
		testGetStatus,
		testGetDriverJourneys,
	}

	for _, f := range testCases {
		testThrowErrorOnTest(t, f)
	}
}
