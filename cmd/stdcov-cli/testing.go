package main

import (
	"io"
	"net/http"
	"strings"

	"gitlab.com/multi/stdcov-api-test/cmd/stdcov-cli/client"
)

type mockClient struct {
	Response      *http.Response
	Error         error
	nCalls        int
	lastURLCalled string
}

func returnErrorClient(err error) APIClient {
	m := &mockClient{Error: err}
	return newTestClient(m)
}

func newTestClient(m *mockClient) *client.Client {
	c, _ := client.NewClient("https://localhost:8000", client.WithHTTPClient(m))
	return c
}

// Get returns the stored response of the mockClient, implements
// HTTPRequestDoer
func (m *mockClient) Do(req *http.Request) (*http.Response, error) {
	m.nCalls++
	m.lastURLCalled = req.URL.String()
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

type NoOpAssertion struct{}

func (NoOpAssertion) Execute() error {
	return nil
}

func (NoOpAssertion) Describe() string {
	return "No assertion"
}
