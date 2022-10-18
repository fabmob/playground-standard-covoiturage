package main

import (
	"context"
	"io"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers/gorillamux"
)

//go:embed ../../stdcov_openapi.yaml
var OpenAPISpec []byte

// ValidateResponse validates a Response against the openapi specification.
func ValidateResponse(request *http.Request, response *http.Response) error {
	ctx := context.Background()
	loader := &openapi3.Loader{Context: ctx, IsExternalRefsAllowed: true}
	spec, loadingErr := loader.LoadFromData(OpenAPISpec)
	panicIf(loadingErr) // Error only if problem with module internals

	specValidationErr := spec.Validate(ctx)
	panicIf(specValidationErr) // Error only if problem with module internals

	router, routerErr := gorillamux.NewRouter(spec)
	panicIf(routerErr) // Error only if problem with module internals

	// Find route
	route, pathParams, err := router.FindRoute(request)
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
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	responseValidationInput.SetBodyBytes(body)
	validationErr := openapi3filter.ValidateResponse(ctx, responseValidationInput)
	return validationErr
}
