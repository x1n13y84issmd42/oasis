package api

import "net/http"

// Response describes a generic HTTP response.
type Response struct {
	StatusCode  int
	ContentType string
	Headers     http.Header
}
