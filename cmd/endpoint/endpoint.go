// Package endpoint provide functions and types to manipulate endpoint and
// server information, to extract this information from a request and to store
// this information into a request's context.
//
// It also stores exported variables for the endpoints of the standard
// covoiturage.
package endpoint

// Info describes an Endpoint
type Info struct {
	Method       string
	Path         string
	HasPathParam bool
}

type Server string

// String implements the Stringer interface for endpoint.Info type
func (e Info) String() string {
	return e.Method + " " + e.Path
}

func New(method, path string) Info {
	return Info{method, path, false}
}

func NewWithParam(method, path string) Info {
	return Info{method, path, true}
}
