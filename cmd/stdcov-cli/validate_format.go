package main

import (
	"context"
	"io"
	"net/http"
	"path/filepath"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers/gorillamux"
)

// OpenAPIv3SpecPath is the relative access path to the OpenAPI specification
var OpenAPIv3SpecPath = filepath.Join("..", "..", "stdcov_openapi.yaml")

// ValidateResponse validates a Response against the openapi specification
func ValidateResponse(request *http.Request, response *http.Response) error {
	ctx := context.Background()
	loader := &openapi3.Loader{Context: ctx, IsExternalRefsAllowed: true}
	doc, err := loader.LoadFromFile(OpenAPIv3SpecPath)

	if err != nil {
		panic(err)
	}

	router, err := gorillamux.NewRouter(doc)

	if err != nil {
		panic(err)
	}

	// Find route
	route, pathParams, err := router.FindRoute(request)

	if err != nil {
		panic(err)
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
		panic(err)
	}

	responseValidationInput.SetBodyBytes(body)
	validationErr := openapi3filter.ValidateResponse(ctx, responseValidationInput)
	return validationErr
}
