package main

import (
	"fmt"
	"os"

	"github.com/labstack/echo/v4"
	"gitlab.com/multi/stdcov-api-test/cmd/stdcov-service/server"
)

//go:generate oapi-codegen -package server -o ./server/server.go -generate "types,server" --old-config-style ../../spec/stdcov_openapi.yaml

func main() {
	handler, err := NewDefaultServer()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	e := echo.New()
	server.RegisterHandlers(e, handler)
	e.Logger.Fatal(e.Start(":1323"))
}
