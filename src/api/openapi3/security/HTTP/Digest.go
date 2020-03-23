package http

import (
	"fmt"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/x1n13y84issmd42/oasis/src/log"
)

// Digest implements a Digest http authentication.
type Digest struct {
	APISecurity *openapi3.SecurityScheme
	Log         log.ILogger
	Auth        WWWAuthenticate
}

// Secure adds an example value from the API spec to the Authorization request header.
func (sec Digest) Secure(req *http.Request) {
	example := sec.APISecurity.Extensions["x-example"].(string)
	if example == "" {
		fmt.Printf("The security \"%s\" contains no example to use in request.", sec.APISecurity)
	} else {
		req.Header["Authorization"] = append(req.Header["Authorization"], example)
	}
}
