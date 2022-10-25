package stdcovcli

import (
	"fmt"
	"net/http"
	"strings"
)

// Query implements flag.Value interface to store query parameters
type Query struct {
	params map[string]string
}

// String implements flag.Value.String
func (qp *Query) String() string {
	str := ""
	for k, v := range qp.params {
		str += fmt.Sprintf("--%s:%s ", k, v)
	}
	return str
}

// Set implements flag.Value.Set
func (qp *Query) Set(s string) error {
	parts := strings.SplitN(s, "=", 2)
	key := parts[0]
	value := ""
	if len(parts) > 1 {
		value = parts[1]
	}
	if qp.params == nil {
		qp.params = make(map[string]string)
	}
	qp.params[key] = value
	return nil
}

// AddQueryParameters adds query parameters stored in a Query object to an
// existing request
func AddQueryParameters(query Query, req *http.Request) {
	q := req.URL.Query()
	for k, v := range query.params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()
}
