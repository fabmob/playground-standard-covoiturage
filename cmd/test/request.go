package test

import (
	"fmt"
	"net/http"
	"strings"
)

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
