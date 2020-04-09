package apikey

import (
	"fmt"
	"net/http"
)

// Query represents an API token security which expects a token as a query parameter.
type Query struct {
	Security
}

// Secure adds an API key to the request's query.
func (sec Query) Secure(req *http.Request) {
	sec.Log.UsingSecurity(sec)

	if sec.Value != "" {
		q := req.URL.Query()
		q.Add(sec.ParamName, sec.Value)
		req.URL.RawQuery = q.Encode()
	} else {
		fmt.Printf("The security \"%s\" contains no data to use in request.", sec.Name)
	}
}
