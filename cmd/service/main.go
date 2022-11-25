// Package service serves a fake API complying with standard covoiturage
// specification.
//
// The server is launched wih the `Run` function, which optionally accepts the
// path to a data file (json format). See Package db documentation for more
// information about the data format. If an empty path is provided, then
// default data is loaded.
package service

import (
	"os"

	"github.com/fabmob/playground-standard-covoiturage/cmd/api"
	"github.com/fabmob/playground-standard-covoiturage/cmd/service/db"
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

		mockDB, err := db.NewMockDBWithData(fileReader)
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
