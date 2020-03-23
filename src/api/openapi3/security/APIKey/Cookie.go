package apikey

import (
	"fmt"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/x1n13y84issmd42/oasis/src/log"
)

// Cookie represents an API token security which expects a token as a header value.
type Cookie struct {
	Log         log.ILogger
	APISecurity openapi3.SecurityScheme
}

// Secure adds an API key to the request's headers.
func (sec Cookie) Secure(req *http.Request) {
	example := sec.APISecurity.Extensions["x-example"].(string)
	if example == "" {
		fmt.Printf("The security \"%s\" contains no example to use in request.", sec.APISecurity)
	} else {
		req.AddCookie(&http.Cookie{Name: sec.APISecurity.Name, Value: example})
	}
}
