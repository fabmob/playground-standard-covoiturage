package main

// Endpoint describes an Endpoint
type Endpoint struct {
	method string
	path   string
}

// String implements the Stringer interface for Endpoint type
func (e Endpoint) String() string {
	return e.method + " " + e.path
}

// panicIf panics if err is not nil
func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}
