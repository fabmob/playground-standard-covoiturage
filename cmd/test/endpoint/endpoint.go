package endpoint

// Info describes an Endpoint
type Info struct {
	Method       string
	Path         string
	HasPathParam bool
}

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
