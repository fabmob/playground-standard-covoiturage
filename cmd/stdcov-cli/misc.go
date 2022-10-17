package main

// Endpoint describes an Endpoint
type Endpoint struct {
	path   string
	method string
}

// String implements the Stringer interface for Endpoint type
func (e Endpoint) String() string {
	return e.method + " " + e.path
}

func panicIfError(err error) {
	if err != nil {
		panic(err)
	}
}
