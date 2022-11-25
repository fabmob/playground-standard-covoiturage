package assert

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/fabmob/playground-standard-covoiturage/cmd/util"
)

func getResponseStatus(obj json.RawMessage) (string, error) {
	type WithStatus struct {
		Status string `json:"status"`
	}

	var withStatus WithStatus

	err := json.Unmarshal(obj, &withStatus)
	if err != nil {
		return "", err
	}

	return withStatus.Status, nil
}

// getQueryCoord extracts departure or arrival coordinates from
// queryParameters
func getQueryCoord(departureOrArrival departureOrArrival, request *http.Request) (util.Coord, error) {
	var latParam, lonParam string

	switch departureOrArrival {
	case departure:
		latParam = "departureLat"
		lonParam = "departureLng"

	case arrival:
		latParam = "arrivalLat"
		lonParam = "arrivalLng"
	}

	lat, err := parseQueryFloatParam(request, latParam)
	if err != nil {
		return util.Coord{}, err
	}

	lon, err := parseQueryFloatParam(request, lonParam)
	if err != nil {
		return util.Coord{}, err
	}

	return util.Coord{Lat: lat, Lon: lon}, nil
}

func getQueryRadius(departureOrArrival departureOrArrival, req *http.Request) (float64, error) {
	const DefaultRadius float64 = 1
	return parseQueryFloatParamWithDefault(req, string(departureOrArrival), DefaultRadius)
}

// getQueryTimeDelta extracts timeDelta parameter from request
func getQueryTimeDelta(req *http.Request) (int, error) {
	const DefaultTimeDelta int = 900
	return parseQueryIntParamWithDefault(req, "timeDelta", DefaultTimeDelta)
}

func getQueryDeparturDate(req *http.Request) (int, error) {
	return parseQueryIntParam(req, "departureDate")
}

func getQueryCount(req *http.Request) (int, error) {
	return parseQueryIntParamWithDefault(req, "count", -1)
}

func parseQueryFloatParam(request *http.Request, paramName string) (float64, error) {
	paramStr := request.URL.Query().Get(paramName)
	return auxParseFloat(paramStr)
}

func parseQueryFloatParamWithDefault(request *http.Request, paramName string, defaultValue float64) (float64, error) {
	paramStr := request.URL.Query().Get(paramName)
	return withDefaultFloat(auxParseFloat)(paramStr, defaultValue)
}

func parseQueryIntParam(request *http.Request, paramName string) (int, error) {
	paramStr := request.URL.Query().Get(paramName)
	return auxParseInt(paramStr)
}

func parseQueryIntParamWithDefault(request *http.Request, paramName string, defaultValue int) (int, error) {
	paramStr := request.URL.Query().Get(paramName)
	return withDefaultInt(auxParseInt)(paramStr, defaultValue)
}

func auxParseFloat(paramStr string) (float64, error) {
	param, err := strconv.ParseFloat(paramStr, 64)
	if err != nil {
		return 0, fmt.Errorf(
			"%s could not be properly parsed as float in query (%w)",
			paramStr,
			err,
		)
	}

	return param, nil
}

func withDefaultFloat(parser func(string) (float64,
	error)) func(string, float64) (float64, error) {

	return func(paramStr string, defaultValue float64) (float64, error) {
		if paramStr == "" {
			return defaultValue, nil
		}

		return parser(paramStr)
	}
}

func auxParseInt(paramStr string) (int, error) {
	param, err := strconv.Atoi(paramStr)
	if err != nil {
		return 0, fmt.Errorf(
			"%s could not be properly parsed as int in query (%w)",
			paramStr,
			err,
		)
	}

	return param, nil
}

func withDefaultInt(parser func(string) (int,
	error)) func(string, int) (int, error) {

	return func(paramStr string, defaultValue int) (int, error) {
		if paramStr == "" {
			return defaultValue, nil
		}

		return parser(paramStr)
	}
}

// failedParsing wraps a parsing error with additional details
func failedParsing(responseOrRequest string, err error) error {
	return fmt.Errorf(
		"internal error while parsing %s:%w",
		responseOrRequest,
		err,
	)
}

// parseArrayResponse parses an array of any type, keeping array elements as
// json.RawMessage
func parseArrayResponse(rsp *http.Response) ([]json.RawMessage, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)

	defer func() { _ = rsp.Body.Close() }()

	if err != nil {
		return nil, err
	}

	var dest []json.RawMessage

	if err := json.Unmarshal(bodyBytes, &dest); err != nil {
		return nil, err
	}

	return dest, nil
}

// getResponseCoord extracts departure or arrival coordinates from
// a json.RawMessage, e.g. as returned by parseArrayResponse. Fails if response has no such coordinates.
func getResponseCoord(departureOrArrival departureOrArrival, obj json.RawMessage) (util.Coord, error) {
	var coordResponse util.Coord

	switch departureOrArrival {
	case departure:
		type WithPassengerPickupCoord struct {
			PassengerPickupLat float64 `json:"passengerPickupLat"`
			PassengerPickupLng float64 `json:"passengerPickupLng"`
		}

		var withPassengerPickupCoord WithPassengerPickupCoord

		err := json.Unmarshal(obj, &withPassengerPickupCoord)
		if err != nil {
			return util.Coord{}, err
		}

		coordResponse = util.Coord{
			Lat: withPassengerPickupCoord.PassengerPickupLat,
			Lon: withPassengerPickupCoord.PassengerPickupLng,
		}

	case arrival:
		type WithPassengerDropCoord struct {
			PassengerDropLat float64 `json:"passengerDropLat"`
			PassengerDropLng float64 `json:"passengerDropLng"`
		}

		var withPassengerDropCoord WithPassengerDropCoord

		err := json.Unmarshal(obj, &withPassengerDropCoord)
		if err != nil {
			return util.Coord{}, err
		}

		coordResponse = util.Coord{
			Lat: withPassengerDropCoord.PassengerDropLat,
			Lon: withPassengerDropCoord.PassengerDropLng,
		}
	}

	return coordResponse, nil
}

func getResponsePickupDate(obj json.RawMessage) (int, error) {
	type WithPickupDate struct {
		PassengerPickupDate int `json:"passengerPickupDate"`
	}

	var withPickupDate WithPickupDate

	err := json.Unmarshal(obj, &withPickupDate)
	if err != nil {
		return 0, err
	}

	return withPickupDate.PassengerPickupDate, nil
}

func getResponseID(obj json.RawMessage) (*string, error) {
	type WithID struct {
		ID *string `json:"id,omitempty"`
	}

	var withID WithID

	err := json.Unmarshal(obj, &withID)
	if err != nil {
		return nil, err
	}

	return withID.ID, nil
}

func getResponseOperator(obj json.RawMessage) (string, error) {
	type WithOperator struct {
		Operator string `json:"operator"`
	}

	var withOperator WithOperator

	err := json.Unmarshal(obj, &withOperator)
	if err != nil {
		return "", err
	}

	return withOperator.Operator, nil
}
