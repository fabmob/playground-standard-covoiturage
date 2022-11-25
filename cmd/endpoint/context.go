package endpoint

import (
	"context"
	"errors"
	"net/http"
)

type key int

const (
	serverKey key = iota
	endpointKey
)

// FromRequest guesses server and endpoint information from a request.
func FromRequest(req *http.Request) (Server, Info, error) {
	return splitURL(req.Method, req.URL.String())
}

// NewContext creates a new context with additional server and endpoint
// information.
func NewContext(ctx context.Context, server Server, endpoint Info) context.Context {
	ctx = context.WithValue(ctx, serverKey, server)
	ctx = context.WithValue(ctx, endpointKey, endpoint)

	return ctx
}

// FromContext retrieves server and endpoint information from request context,
// or returns an error if missing.
func FromContext(ctx context.Context) (Server, Info, error) {
	server, okServer := ctx.Value(serverKey).(Server)
	endpoint, okEndpoint := ctx.Value(endpointKey).(Info)

	var err error

	if !okServer || !okEndpoint {
		err = errors.New("could not get endpoint and server from request context")
	}

	return server, endpoint, err
}
