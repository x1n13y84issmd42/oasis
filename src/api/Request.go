package api

import (
	"net/http"
)

// Request describes the properties of an HTTP request.
type Request struct {
	Method  string
	Path    string
	Headers http.Header
}
