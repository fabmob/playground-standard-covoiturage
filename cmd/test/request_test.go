package test

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

var URL = localServer + "/status"

func TestMakeRequestXAPIKey(t *testing.T) {
	var (
		method        = http.MethodGet
		body   []byte = nil
	)

	testCases := []string{
		"1234",
		"4567",
	}

	for _, apiKey := range testCases {
		req, err := makeRequest(method, URL, body, apiKey)
		panicIf(err)

		if req.Header.Get("X-API-Key") != apiKey {
			t.Error("X-API-Key header is not specified properly")
		}
	}

}

func TestMakeRequestBody(t *testing.T) {
	bodyStr := "test body"
	bodyBytes := []byte(bodyStr)

	req, err := makeRequest(http.MethodGet, URL, bodyBytes, "")
	panicIf(err)

	if req.Body == nil {
		t.Fatal("makeRequest does not initializes the body properly")
	}

	body, err := io.ReadAll(req.Body)

	if err != nil || string(body) != bodyStr {
		t.Error("makeRequest does not initializes the body properly")
	}
}

func TestMakeRequestHeader(t *testing.T) {
	bodyStr := "test body"
	bodyBytes := []byte(bodyStr)

	req, err := makeRequest(http.MethodGet, URL, bodyBytes, "")
	panicIf(err)

	if !strings.HasPrefix(req.Header.Get(echo.HeaderContentType), echo.MIMEApplicationJSON) {
		t.Fail()
	}

}
