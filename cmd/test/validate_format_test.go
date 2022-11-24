package test

import (
	"context"
	"net/http"
	"testing"

	"github.com/fabmob/playground-standard-covoiturage/cmd/api"
)

func TestUndocumentedStatusCode(t *testing.T) {
	var (
		invalidStatus = http.StatusLoopDetected
	)
	request, err := api.NewGetDriverJourneysRequest(localServer, &api.GetDriverJourneysParams{})
	panicIf(err)

	response := mockResponse(invalidStatus, "[]", nil)

	validationErr := validateResponse(request, response)
	if validationErr == nil {
		t.Error("Format validation is expected to fail for undocumented status code")
	}

}

func TestFindRoute(t *testing.T) {
	server := "https://abc.fr/abc"
	url := server + "/driver_journeys"

	request, err := http.NewRequest("GET", url, nil)
	panicIf(err)

	ctx := context.Background()

	_, _, err = findRoute(ctx, request, server)
	if err != nil {
		t.Log(server)
		t.Error("Format validation with kin-openapi does not find route properly")
	}

}
