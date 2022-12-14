package test

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/fabmob/playground-standard-covoiturage/cmd/test/assert"
	"github.com/fabmob/playground-standard-covoiturage/cmd/util"
	"github.com/google/go-cmp/cmp"
)

var defaultTestFlags Flags = Flags{ExpectNonEmpty: false}

// testErrorOnRequestIsHandled returns an urlError for every API call and checks:
// - that only one AssertionError is returned
// - that AssertionError.Unwrap() != nil
func testErrorOnRequestIsHandled(t *testing.T, f requestTestFun) {
	t.Helper()

	t.Run("API call throws error", func(t *testing.T) {
		urlError := &url.Error{Op: "", URL: "", Err: errors.New("error")}
		m := NewMockClientWithError(urlError)

		// specific request is irrelevant as the error client is in any case returning an error
		r, err := http.NewRequest(http.MethodGet, "/", strings.NewReader(""))
		util.PanicIf(err)

		results := f(m, r, defaultTestFlags)
		assert.ShouldHaveSingleAssertionResult(t, results)

		err = results[0].Unwrap()
		if err == nil {
			t.Error("If error returned, api is not up")
		}
	})
}

func TestAPIErrors(t *testing.T) {
	for _, fun := range GetAPIMapping() {
		testErrorOnRequestIsHandled(t, wrapTestResponseFun(fun))
	}
}

// Test that the expected requests are made with wrapTest
func TestRequests(t *testing.T) {
	testCases := []string{
		"/driver_journeys",
		"/driver_journeys?departureLat=0&departureLng=0",
	}

	for _, url := range testCases {
		t.Run(url, func(t *testing.T) {
			m := NewMockClientWithResponse(assert.MockOKStatusResponse())
			r, err := http.NewRequest(http.MethodGet, url, strings.NewReader(""))
			util.PanicIf(err)

			testNoAssertions := func(*http.Request, *http.Response, Flags) []assert.Result {
				return nil
			}
			wrapTestResponseFun(testNoAssertions)(m, r, defaultTestFlags)

			requestsDone := m.Client.(*MockClient).Requests
			if len(requestsDone) != 1 {
				t.Error("MockClient is expected to do exactly one request")
			}

			if !cmpRequests(t, requestsDone[0], r) {
				t.Logf("Request done: %+v", requestsDone[0])
				t.Logf("Request expected: %+v", r)
				t.Error("MockClient request is expected to be the one passed as argument")
			}
		})
	}
}

// cmpRequests ensures req1 and req2 have same Method, url and headers.
func cmpRequests(t *testing.T, req1, req2 *http.Request) bool {
	t.Helper()

	body := make([]io.Reader, 2)
	bodyString := make([]string, 2)

	reqs := []*http.Request{req1, req2}
	for i, req := range reqs {
		var err error

		body[i], err = req.GetBody()
		util.PanicIf(err)

		bodyBytes, err := io.ReadAll(body[i])
		util.PanicIf(err)

		bodyString[i] = string(bodyBytes)
	}

	return req1.Method == req2.Method &&
		req1.URL.String() == req2.URL.String() &&
		cmp.Equal(req1.Header, req2.Header) &&
		bodyString[0] == bodyString[1]
}
