package main

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

//  testErrorOnRequestIsHandled returns an urlError for every API call and checks:
// - that only one AssertionError is returned
// - that AssertionError.Unwrap() != nil
func testErrorOnRequestIsHandled(t *testing.T, f auxTestFun) {
	t.Helper()
	t.Run("API call throws error", func(t *testing.T) {
		urlError := &url.Error{Op: "", URL: "", Err: errors.New("error")}
		m := NewMockClientWithError(urlError)
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
		testErrorOnRequestIsHandled(t, f)
	}
}

// Test that the expected requests are made.
func TestRequests(t *testing.T) {
	m := NewMockClientWithResponse(mockOKStatusResponse())
	r, err := http.NewRequest(http.MethodGet, "/driver_journeys", strings.NewReader(""))
	panicIf(err)
	a := NewAssertionAccu()
	testGetDriverJourneys(m, r, a)

	requestsDone := m.Client.(*MockClient).Requests
	if len(requestsDone) != 1 {
		t.Error("MockClient is expected to do exactly one request")
	}
	if requestsDone[0] != r {
		t.Logf("Request done: %+v", requestsDone[0])
		t.Logf("Request expected: %+v", r)
		t.Error("MockClient request is expected to be the one passed as argument")
	}

}
