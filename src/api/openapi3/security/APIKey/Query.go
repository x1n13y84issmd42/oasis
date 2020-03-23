package apikey

import (
	"fmt"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/x1n13y84issmd42/oasis/src/log"
)

// Query represents an API token security which expects a token as a query parameter.
type Query struct {
	Log         log.ILogger
	APISecurity openapi3.SecurityScheme
}

// Secure adds an API key to the request's query.
func (sec Query) Secure(req *http.Request) {
	example := sec.APISecurity.Extensions["x-example"].(string)
	if example == "" {
		fmt.Printf("The security \"%s\" contains no example to use in request.\n", sec.APISecurity)
	} else {
		q := req.URL.Query()
		q.Add(sec.APISecurity.Name, example)
		req.URL.RawQuery = q.Encode()
	}
}
