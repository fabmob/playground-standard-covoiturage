package main

import (
	"io"
	"net/http"
	"strings"

	"gitlab.com/multi/stdcov-api-test/cmd/stdcov-cli/client"
)

type mockClient struct {
	Response *http.Response
	Error    error
	Requests []*http.Request
}

func NewMockClientWithError(err error) APIClient {
	m := &mockClient{Error: err}
	return newTestClient(m)
}

func NewMockClientWithResponse(r *http.Response) APIClient {
	m := &mockClient{Response: r}
	return newTestClient(m)
}

func newTestClient(m *mockClient) *client.Client {
	c, _ := client.NewClient("", client.WithHTTPClient(m))
	return c
}

// Get returns the stored response of the mockClient, implements
// HTTPRequestDoer
func (m *mockClient) Do(req *http.Request) (*http.Response, error) {
	m.Requests = append(m.Requests, req)
	if m.Error != nil {
		return nil, m.Error
	}
	return m.Response, nil
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

// A NoOpAssertion returns stored error when executed
type NoOpAssertion struct{ error }

// Execute implements Assertion interface
func (n NoOpAssertion) Execute() error {
	return n.error
}

// Describe implements Assertion interface
func (NoOpAssertion) Describe() string {
	return "No assertion"
}
