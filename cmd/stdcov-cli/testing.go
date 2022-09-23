package main

import (
	"net/http"
)

type mockClient struct {
	Response      *http.Response
	Error         error
	nCalls        int
	lastURLCalled string
}

// returnStatusCodeClient is a client that always return an empty response
// with given status code
func returnStatusCodeClient(statusCode int) *mockClient {
	response := mockResponse(statusCode)
	m := &mockClient{Response: response}
	return m
}

// Get returns the stored response of the mockClient
func (m *mockClient) Get(url string) (*http.Response, error) {
	m.nCalls++
	m.lastURLCalled = url
	if m.Error != nil {
		return nil, m.Error
	}
	return m.Response, nil
}

func mockResponse(statusCode int) *http.Response {
	return &http.Response{
		Status:     http.StatusText(statusCode),
		StatusCode: statusCode,
	}
}
