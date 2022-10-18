package main

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/multi/stdcov-api-test/cmd/stdcov-service/server"
)

//go:generate oapi-codegen -package server -o ./server/server.go -generate "types,server" --old-config-style ../../spec/stdcov_openapi.yaml

func main() {
	var api *StdCovServerImpl
	e := echo.New()
	server.RegisterHandlers(e, api)
	e.Logger.Fatal(e.Start(":1323"))
}
