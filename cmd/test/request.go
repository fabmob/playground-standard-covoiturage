package test

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"github.com/fabmob/playground-standard-covoiturage/cmd/endpoint"
	"github.com/labstack/echo/v4"
)

const (
	HeaderXAPIKey = "X-API-Key"
)

// makeRequestWithContext makes a request with right header, and stores server and
// endpoint information in its context.
// Assumes URL is pointing towards a valid endpoint of the standard-covoiturage.
func makeRequestWithContext(method, URL string, body []byte, apiKey string) (*http.Request, error) {
	req, err := http.NewRequest(method, URL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set(HeaderXAPIKey, apiKey)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	return AddEndpointContext(req)
}

// AddEndpointContext adds server and endpoint information to the request's
// context. Use endpoint.FromContext to extract this information.
func AddEndpointContext(request *http.Request) (*http.Request, error) {
	server, endpointInfo, err := endpoint.FromRequest(request)
	if err != nil {
		return nil, err
	}

	ctx := endpoint.NewContext(request.Context(), server, endpointInfo)

	return request.WithContext(ctx), nil
}

// Query implements flag.Value interface to store query parameters
type Query struct {
	Params map[string]string
}

// NewQuery initializes an empty query
func NewQuery() Query {
	return Query{map[string]string{}}
}

// String implements pflag.Value.String (cobra flags)
func (qp *Query) String() string {
	var str = ""

	for k, v := range qp.Params {
		str += fmt.Sprintf("--%s:%s ", k, v)
	}

	return str
}

// SetParam sets a query parameter. If key already exists, it is overwritten.
func (qp *Query) SetParam(key, value string) {
	qp.Params[key] = value
}

// SetOptionalParam sets a query parameter, only if the value is not "". If
// key already exists, it may be overwritten.
func (qp *Query) SetOptionalParam(key, value string) {
	if value != "" {
		qp.Params[key] = value
	}
}

// Set implements pflag.Value.Set (cobra flags)
func (qp *Query) Set(s string) error {
	parts := strings.SplitN(s, "=", 2)
	key := parts[0]
	value := ""

	if len(parts) > 1 {
		value = parts[1]
	}

	if qp.Params == nil {
		qp.Params = make(map[string]string)
	}

	qp.Params[key] = value

	return nil
}

// Type implements pflag.Value.Type (cobra flags)
func (qp *Query) Type() string {
	return "*Query"
}

// AddQueryParameters adds query parameters stored in a Query object to an
// existing request
func AddQueryParameters(query Query, req *http.Request) {
	q := req.URL.Query()

	for k, v := range query.Params {
		q.Add(k, v)
	}

	req.URL.RawQuery = q.Encode()
}
