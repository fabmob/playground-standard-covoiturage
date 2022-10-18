package spec

import (
	_ "embed" // for go:embed directive
)

// OpenAPISpec is the OpenAPI specification
//
//go:embed stdcov_openapi.yaml
var OpenAPISpec []byte
