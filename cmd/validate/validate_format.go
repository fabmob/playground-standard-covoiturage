package validate

import (
	"context"
	"io"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers/gorillamux"

	"gitlab.com/multi/stdcov-api-test/spec"
)

// Response validates a Response against the openapi specification.
func Response(request *http.Request, response *http.Response) error {
	ctx := context.Background()
	loader := &openapi3.Loader{Context: ctx, IsExternalRefsAllowed: true}
	doc, loadingErr := loader.LoadFromData(spec.OpenAPISpec)
	if loadingErr != nil {
		panic(loadingErr) // Error only if problem with module internals
	}

	specValidationErr := doc.Validate(ctx)
	if specValidationErr != nil {
		panic(specValidationErr) // Error only if problem with module internals
	}

	router, routerErr := gorillamux.NewRouter(doc)
	if routerErr != nil {
		panic(routerErr) // Error only if problem with module internals
	}

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
