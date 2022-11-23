package service

import (
	"encoding/json"
	"os"

	"github.com/fabmob/playground-standard-covoiturage/cmd/api"
	"github.com/labstack/echo/v4"
)

// Run serves a server with an implementation of the API enforcing the
// "standard-covoiturage" specification
func Run(dataFile string) {
	e := echo.New()

	var handler *StdCovServerImpl

	if dataFile == "" {
		handler = NewDefaultServer()
	} else {

		fileBytes, err := os.ReadFile(dataFile)
		if err != nil {
			e.Logger.Fatal(err)
		}

		var mockDB MockDB

		json.Unmarshal(fileBytes, &mockDB)

		handler = &StdCovServerImpl{&mockDB}
	}

	api.RegisterHandlers(e, handler)
	e.Logger.Fatal(e.Start(":1323"))
}
