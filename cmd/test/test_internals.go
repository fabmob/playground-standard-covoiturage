package test

import (
	"net/http"

	"gitlab.com/multi/stdcov-api-test/cmd/api"
)

// APIClient is a client to the API standard covoiturage
type APIClient = *api.Client

// Request tests a request
func Request(client APIClient, request *http.Request, flags Flags) (*Report, error) {
	endpoint, err := ExtractEndpoint(request, client.Server)
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
	all := []AssertionResult{}
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
		a := NewAssertionAccu()
		f(req, resp, a, flags)
		a.ExecuteAll()
		return a.GetAssertionResults()
	}
}