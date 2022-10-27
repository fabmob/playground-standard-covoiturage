package service

import (
	"fmt"
	"os"

	"github.com/labstack/echo/v4"
	"gitlab.com/multi/stdcov-api-test/cmd/api"
)

func Run() {
	handler, err := NewDefaultServer()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	e := echo.New()
	api.RegisterHandlers(e, handler)
	e.Logger.Fatal(e.Start(":1323"))
}
