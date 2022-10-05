package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
)

func main() {
	urlStrPtr := flag.String("url", "", "Base url of the API under test")

	flag.Parse()

	baseURL, _ := url.Parse(*urlStrPtr)
	client := &http.Client{}
	fmt.Println(Check(client, baseURL))
}

// An HTTPGetter can issue HTTP Get requests
type HTTPGetter interface {
	Get(url string) (*http.Response, error)
}
