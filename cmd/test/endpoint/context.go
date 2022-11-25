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

func FromRequest(req *http.Request) (Server, Info, error) {
	return SplitURL(req.Method, req.URL.String())
}

func NewContext(ctx context.Context, server Server, endpoint Info) context.Context {
	ctx = context.WithValue(ctx, serverKey, server)
	ctx = context.WithValue(ctx, endpointKey, endpoint)

	return ctx
}

func FromContext(ctx context.Context) (Server, Info, error) {
	server, okServer := ctx.Value(serverKey).(Server)
	endpoint, okEndpoint := ctx.Value(endpointKey).(Info)

	var err error

	if !okServer || !okEndpoint {
		err = errors.New("could not get endpoint and server from request context")
	}

	return server, endpoint, err
}
