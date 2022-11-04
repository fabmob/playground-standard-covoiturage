package test

import (
	"net/http"

	"github.com/fabmob/playground-standard-covoiturage/cmd/api"
)

// APIClient is a client to the API standard covoiturage
type APIClient = *api.Client

// Request tests a request
func Request(server string, request *http.Request, flags Flags) (*Report, error) {

	client, err := api.NewClient(server)
	if err != nil {
		return nil, err
	}

	endpoint, err := ExtractEndpoint(request, server)
	if err != nil {
		return nil, err
	}

	selectedTestFun, err := SelectTestFuns(endpoint)
	if err != nil {
		return nil, err
	}

	report := executeTestFuns(client, request, selectedTestFun, flags)
	report.endpoint = endpoint

	return report, nil
}

func executeTestFuns(
	client APIClient,
	request *http.Request,
	testFun ResponseTestFun,
	flags Flags,
) *Report {
	var all = []AssertionResult{}

	all = append(all, wrapTestResponseFun(testFun)(client, request, flags)...)
	report := NewReport(all...)

	return &report
}

/////////////////////////////////////////////////////////////

// A requestTestFun runs all tests associated with a given Request, and return
// the correspending `AssertionResult`s
type requestTestFun func(APIClient, *http.Request, Flags) []AssertionResult

// wrapTestResponseFun wraps an TestResponseFun to a TestRequestFun
func wrapTestResponseFun(f ResponseTestFun) requestTestFun {
	return func(c APIClient, request *http.Request, flags Flags) []AssertionResult {
		response, clientErr := c.Client.Do(request)
		if clientErr != nil {
			return []AssertionResult{CheckAPICallSuccess(clientErr)}
		}

		return f(request, response, flags)
	}
}

//////////////////////////////////////////////////////////////

// A ResponseTestFun runs tests on a given *http.Response (given a
// *http.Request) and return the correspending `AssertionsResult`s
type ResponseTestFun func(
	*http.Request,
	*http.Response,
	Flags,
) []AssertionResult

func wrapAssertionsFun(f testImplementation) ResponseTestFun {
	return func(req *http.Request, resp *http.Response, flags Flags) []AssertionResult {
		var (
			err error
			a   = NewAssertionAccu()
		)

		resp.Body, err = ReusableReadCloser(resp.Body)
		if err != nil {
			return []AssertionResult{NewAssertionResult(err, "failure to read response")}
		}

		f(req, resp, a, flags)

		a.ExecuteAll()

		return a.GetAssertionResults()
	}
}
