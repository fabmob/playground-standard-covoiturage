package service

import (
	"github.com/fabmob/playground-standard-covoiturage/cmd/api"
	"github.com/labstack/echo/v4"
)

// Run serves a server with an implementation of the API enforcing the
// "standard-covoiturage" specification
func Run() {
	handler := NewDefaultServer()
	e := echo.New()
	api.RegisterHandlers(e, handler)
	e.Logger.Fatal(e.Start(":1323"))
}
