package test

import (
	"net/http"

	"github.com/fabmob/playground-standard-covoiturage/cmd/api"
	"github.com/fabmob/playground-standard-covoiturage/cmd/endpoint"
	"github.com/fabmob/playground-standard-covoiturage/cmd/test/assert"
)

// APIClient is a client to the API standard covoiturage
type APIClient = *api.Client

// testRequest tests a testRequest
func testRequest(request *http.Request, flags Flags) (*Report, error) {

	server, endpoint, err := endpoint.FromContext(request.Context())
	if err != nil {
		return nil, err
	}

	client, err := api.NewClient(string(server))
	if err != nil {
		return nil, err
	}

	selectedTestFun, err := SelectTestFun(endpoint)
	if err != nil {
		return nil, err
	}

	report := executeTestFun(client, request, selectedTestFun, flags)

	report.endpoint = endpoint

	return report, nil
}

func executeTestFun(
	client APIClient,
	request *http.Request,
	testFun ResponseTestFun,
	flags Flags,
) *Report {
	var all = []assert.Result{}

	all = append(all, wrapTestResponseFun(testFun)(client, request, flags)...)
	report := NewReport(all...)

	return &report
}

/////////////////////////////////////////////////////////////

// A requestTestFun runs all tests associated with a given Request, and
// returns the correspending `assert.Result`s
type requestTestFun func(APIClient, *http.Request, Flags) []assert.Result

// wrapTestResponseFun wraps a TestResponseFun to a TestRequestFun
func wrapTestResponseFun(f ResponseTestFun) requestTestFun {

	return func(c APIClient, request *http.Request, flags Flags) []assert.Result {
		response, clientErr := c.Client.Do(request)
		if clientErr != nil {
			return []assert.Result{assert.CheckAPICallSuccess(clientErr)}
		}

		return f(request, response, flags)
	}
}

//////////////////////////////////////////////////////////////

// A ResponseTestFun runs tests on a given *http.Response (given a
// *http.Request) and return the correspending `assert.sResult`s
type ResponseTestFun func(
	*http.Request,
	*http.Response,
	Flags,
) []assert.Result

func wrapTestImplementation(f testImplementation) ResponseTestFun {

	return func(req *http.Request, resp *http.Response, flags Flags) []assert.Result {
		var (
			err error
			a   = assert.NewAccumulator()
		)

		// response body may be read several times in assertions
		resp.Body, err = ReusableReadCloser(resp.Body)
		if err != nil {
			return []assert.Result{assert.NewAssertionResult(err, "failure to read response body")}
		}

		f(req, resp, a, flags)

		a.ExecuteAll()

		return a.GetAssertionResults()
	}
}
