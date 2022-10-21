package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"gitlab.com/multi/stdcov-api-test/cmd/stdcov-cli/client"
)

//go:generate oapi-codegen -package client -o ./client/client.go -generate "types,client" --old-config-style ../../spec/stdcov_openapi.yaml

func main() {
	os.Exit(Run())
}

// Run runs the cli validation and returns an exit code
func Run() int {
	urlStrPtr := flag.String("url", "", "Server url of the API under test")
	verboseBoolPtr := flag.Bool("verbose", false, "Make the operation more talkative")
	var query Query
	flag.Var(&query, "q", "Query parameters in the form name=value")

	flag.Parse()

	c, err := client.NewClient("")
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("GET", *urlStrPtr, nil)
	if err != nil {
		fmt.Println(err)
		return 1
	}

	q := req.URL.Query()
	for k, v := range query.params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	report, err := ExecuteTestSuite(c, req)
	if err != nil {
		fmt.Println(err)
		return 1
	}

	report.verbose = *verboseBoolPtr
	fmt.Println(report)
	if report.hasErrors() {
		return 1
	}
	return 0
}

// Query implements flag.Value interface to store query parameters
type Query struct {
	params map[string]string
}

func (qp *Query) String() string {
	str := ""
	for k, v := range qp.params {
		str += fmt.Sprintf("--%s:%s ", k, v)
	}
	return str
}

func (qp *Query) Set(s string) error {
	parts := strings.SplitN(s, "=", 2)
	key := parts[0]
	value := ""
	if len(parts) > 1 {
		value = parts[1]
	}
	if qp.params == nil {
		qp.params = make(map[string]string)
	}
	qp.params[key] = value
	return nil
}
