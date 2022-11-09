package test

import (
	"net/http"
	"testing"

	"github.com/fabmob/playground-standard-covoiturage/cmd/api"
)

func TestWrongStatusCode(t *testing.T) {
	var (
		fakeServer    = ""
		invalidStatus = http.StatusLoopDetected
	)
	request, err := api.NewGetDriverJourneysRequest(fakeServer, &api.GetDriverJourneysParams{})
	panicIf(err)

	response := mockResponse(invalidStatus, "[]", nil)

	validationErr := validateResponse(request, response)
	if validationErr == nil {
		t.Error("Format validation is expected to fail for undocumented status code")
	}

}
