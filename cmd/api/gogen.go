// Package api has helper functions to manipulate operations linked to the api
// specification.
//
// Most of the functions are generated from the OpenAPI specfication with
// oapi-codegen.
//
// The package also exports:
//
// * helpers to perform conversions between
// Driver/PassengerCarpoolBooking and Booking types.
// * umbrella type `Journeys` to manipulate driverJourneys and passengerJourneys alike
package api

//go:generate oapi-codegen -config oapi-codegen-config.yaml ../../spec/stdcov_openapi.yaml
