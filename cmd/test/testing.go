package test

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"gitlab.com/multi/stdcov-api-test/cmd/api"
)

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

// mockResponse returns a mock response with given statusCode, body, and
// header. If headers are `nil` default headers with "Content-Type: json" are
// used.
func mockResponse(
	statusCode int,
	body string,
	header http.Header,
) *http.Response {

	if header == nil {
		header = make(http.Header)
		header["Content-Type"] = []string{"json"}
	}

	return &http.Response{
		Status:        http.StatusText(statusCode),
		StatusCode:    statusCode,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Body:          io.NopCloser(strings.NewReader(body)),
		ContentLength: int64(len(body)),
		Header:        header,
	}
}

func mockStatusResponse(statusCode int) *http.Response {
	return mockResponse(statusCode, "", nil)
}

func mockOKStatusResponse() *http.Response {
	return mockStatusResponse(http.StatusOK)
}

func mockBodyResponse(responseObj interface{}) *http.Response {
	responseJSON, err := json.Marshal(responseObj)
	panicIf(err)
	return mockResponse(200, string(responseJSON), nil)
}

// A NopAssertion returns stored error when executed
type NopAssertion struct{ error }

// Execute implements Assertion interface
func (n NopAssertion) Execute() error {
	return n.error
}

// Describe implements Assertion interface
func (NopAssertion) Describe() string {
	return "No assertion"
}

// panicIf panics if err is not nil
func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}
