package contract

import "net/http"

// Expectation is a function to check various properties of an HTTP request.
// A list of various expectation compose a Validator.
type Expectation func(*http.Response) bool
