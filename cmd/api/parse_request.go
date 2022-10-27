package api

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/labstack/echo/v4"
)

// Code is mostly copy-pasted from generated cmd/stdcov-service/server/server.go
// - Functions are renamed to Parse...Request
// - Functions inputs are changed to http.Request (use of FakeContext to mimic
// echo.context)
// - Functions outputs are changed to *...Params of interest
// and all subsequent changes in order to build

// FakeContext mimics echo.Context to change less code by hand
type FakeContext struct {
	request *http.Request
}

// QueryParams returns query parameters
func (fc FakeContext) QueryParams() url.Values {
	return fc.request.URL.Query()
}

// ParseGetDriverJourneysRequest converts an *http.Request into a
// GetDriverJourneysParams
func ParseGetDriverJourneysRequest(req *http.Request) (*GetDriverJourneysParams, error) {
	var err error

	// Parameter object where we will unmarshal all parameters from the request
	var params GetDriverJourneysParams
	ctx := FakeContext{req}
	// ------------- Required query parameter "departureLat" -------------

	err = runtime.BindQueryParameter("form", true, true, "departureLat", ctx.QueryParams(), &params.DepartureLat)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter departureLat: %s", err))
	}

	// ------------- Required query parameter "departureLng" -------------

	err = runtime.BindQueryParameter("form", true, true, "departureLng", ctx.QueryParams(), &params.DepartureLng)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter departureLng: %s", err))
	}

	// ------------- Required query parameter "arrivalLat" -------------

	err = runtime.BindQueryParameter("form", true, true, "arrivalLat", ctx.QueryParams(), &params.ArrivalLat)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter arrivalLat: %s", err))
	}

	// ------------- Required query parameter "arrivalLng" -------------

	err = runtime.BindQueryParameter("form", true, true, "arrivalLng", ctx.QueryParams(), &params.ArrivalLng)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter arrivalLng: %s", err))
	}

	// ------------- Required query parameter "departureDate" -------------

	err = runtime.BindQueryParameter("form", true, true, "departureDate", ctx.QueryParams(), &params.DepartureDate)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter departureDate: %s", err))
	}

	// ------------- Optional query parameter "timeDelta" -------------

	err = runtime.BindQueryParameter("form", true, false, "timeDelta", ctx.QueryParams(), &params.TimeDelta)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter timeDelta: %s", err))
	}

	// ------------- Optional query parameter "departureRadius" -------------

	err = runtime.BindQueryParameter("form", true, false, "departureRadius", ctx.QueryParams(), &params.DepartureRadius)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter departureRadius: %s", err))
	}

	// ------------- Optional query parameter "arrivalRadius" -------------

	err = runtime.BindQueryParameter("form", true, false, "arrivalRadius", ctx.QueryParams(), &params.ArrivalRadius)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter arrivalRadius: %s", err))
	}

	// ------------- Optional query parameter "count" -------------

	err = runtime.BindQueryParameter("form", true, false, "count", ctx.QueryParams(), &params.Count)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter count: %s", err))
	}
	return &params, nil
}
