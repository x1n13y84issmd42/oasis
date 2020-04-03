package api

// Response describes a generic HTTP response.
type Response struct {
	Description string
	StatusCode  uint64
	ContentType string
	Headers     Headers
	Schema      *Schema
}
