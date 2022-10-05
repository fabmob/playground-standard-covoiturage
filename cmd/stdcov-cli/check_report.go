package main

import (
	"net/http"
	"net/url"
)

// Check tests a client against all implemented tests
func Check(client *http.Client, baseURL *url.URL) CheckReport {
	return CheckReport{
		statusCheck: CheckAPIStatus(client, baseURL),
	}
}

// CheckReport stores check results
type CheckReport struct {
	statusCheck bool
}

func (cr CheckReport) String() string {
	if !cr.statusCheck {
		return "❌ Check of the status endpoint failed\n"
	}
	return "✅ Status check OK\n"
}
