package main

import (
	"fmt"
	"net/http"
	"net/url"
)

func main() {
	baseURLStr := "https://json.validator.validata.fr"
	baseURL, _ := url.Parse(baseURLStr)
	client := &http.Client{}
	fmt.Println(CheckAPIStatus(client, baseURL))
}

// An HTTPGetter can issue HTTP Get requests
type HTTPGetter interface {
	Get(url string) (*http.Response, error)
}
