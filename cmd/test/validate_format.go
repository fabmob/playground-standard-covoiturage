package test

import (
	"context"
	"io"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"github.com/getkin/kin-openapi/routers/gorillamux"

	"github.com/fabmob/playground-standard-covoiturage/cmd/endpoint"
	"github.com/fabmob/playground-standard-covoiturage/spec"
)

// Response validates a Response against the openapi specification.
func validateResponse(request *http.Request, response *http.Response) error {
	server, _, err := endpoint.FromContext(request.Context())
	if err != nil {
		return err
	}

	ctx := context.Background()

	route, pathParams, err := findRoute(ctx, request, server)
	if err != nil {
		return err
	}

	requestValidationInput := &openapi3filter.RequestValidationInput{
		Request:    request,
		PathParams: pathParams,
		Route:      route,
	}

	// Validate response
	responseValidationInput := &openapi3filter.ResponseValidationInput{
		RequestValidationInput: requestValidationInput,
		Status:                 response.StatusCode,
		Header:                 response.Header,
		Options:                &openapi3filter.Options{IncludeResponseStatus: true},
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	responseValidationInput.SetBodyBytes(body)
	validationErr := openapi3filter.ValidateResponse(ctx, responseValidationInput)

	return validationErr
}

func findRoute(ctx context.Context, request *http.Request, server endpoint.Server) (route *routers.Route, pathParams map[string]string, err error) {
	loader := &openapi3.Loader{Context: ctx, IsExternalRefsAllowed: true}

	doc, loadingErr := loader.LoadFromData(spec.OpenAPISpec)
	if loadingErr != nil {
		panic(loadingErr) // Error only if problem with module internals
	}

	doc.Servers = openapi3.Servers{&openapi3.Server{URL: string(server)}}

	specValidationErr := doc.Validate(ctx)
	if specValidationErr != nil {
		panic(specValidationErr) // Error only if problem with module internals
	}

	router, routerErr := gorillamux.NewRouter(doc)
	if routerErr != nil {
		panic(routerErr) // Error only if problem with module internals
	}

	return router.FindRoute(request)
}
