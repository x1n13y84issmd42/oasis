package apikey

import (
	"fmt"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/x1n13y84issmd42/oasis/src/log"
)

// Header represents an API token security which expects a token as a header value.
type Header struct {
	Log         log.ILogger
	APISecurity openapi3.SecurityScheme
}

// Secure adds an API key to the request's headers.
func (sec Header) Secure(req *http.Request) {
	example := sec.APISecurity.Extensions["x-example"].(string)
	if example == "" {
		fmt.Printf("The security \"%s\" contains no example to use in request.", sec.APISecurity)
	} else {
		req.Header.Add(sec.APISecurity.Name, example)
	}
}
