package main

import "net/http"

type customTransport struct {
	base    http.RoundTripper
	headers http.Header
}

func (t *customTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	for k, v := range t.headers {
		req.Header[k] = v
	}
	return t.base.RoundTrip(req)
}
