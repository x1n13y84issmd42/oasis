package api

import (
	"net/http"
	"net/url"
)

// Request describes the properties of an HTTP request.
type Request struct {
	Method  string
	URL     url.URL
	Headers http.Header
}
