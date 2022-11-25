package test

import (
	"net/http"

	"github.com/fabmob/playground-standard-covoiturage/cmd/api"
)

const localServer = "http://localhost:1323"

// MockClient is an HTTP client that returns always the same response or
// error, and stores the requests that are made.
type MockClient struct {
	Response *http.Response
	Error    error
	Requests []*http.Request
}

// Do returns the stored response of the MockClient, implements
// HTTPRequestDoer
func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	m.Requests = append(m.Requests, req)

	if m.Error != nil {
		return nil, m.Error
	}

	return m.Response, nil
}

// NewMockClientWithError returns a MockClient that always returns error
// `err`
func NewMockClientWithError(err error) APIClient {
	m := &MockClient{Error: err}
	return newTestClient(m)
}

// NewMockClientWithResponse returns a MockClient that always returns response
// `r`
func NewMockClientWithResponse(r *http.Response) APIClient {
	m := &MockClient{Response: r}
	return newTestClient(m)
}

func newTestClient(m *MockClient) *api.Client {
	c, _ := api.NewClient("", api.WithHTTPClient(m))
	return c
}

//////////////////////////////////////////////////////////////
// Mock Runner
//////////////////////////////////////////////////////////////

// A MockRunner implements TestRunner interface
type MockRunner struct {
	Method  string
	URL     string
	Query   Query
	Body    []byte
	Verbose bool
	APIKey  string
	Flags   Flags
}

// Run stores arguments and returns nil
func (mr *MockRunner) Run(
	method,
	URL string,
	query Query,
	body []byte,
	verbose bool,
	apiKey string,
	flags Flags,
) error {

	mr.Method = method
	mr.URL = URL
	mr.Verbose = verbose
	mr.Query = query
	mr.Body = body
	mr.APIKey = apiKey
	mr.Flags = flags

	return nil
}

func NewMockRunner() *MockRunner {
	return &MockRunner{}
}
