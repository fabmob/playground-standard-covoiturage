package service

import (
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
		fileReader, err := os.Open(dataFile)
		exitIfErr(err, e)

		mockDB, err := NewMockDBWithData(fileReader)
		exitIfErr(err, e)

		handler = NewServerWithDB(mockDB)
	}

	api.RegisterHandlers(e, handler)
	e.Logger.Fatal(e.Start(":1323"))
}

func exitIfErr(err error, e *echo.Echo) {
	if err != nil {
		e.Logger.Fatal(err)
		os.Exit(1)
	}
}
