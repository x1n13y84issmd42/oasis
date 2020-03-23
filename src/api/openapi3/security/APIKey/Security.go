package apikey

import (
	"fmt"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/x1n13y84issmd42/oasis/src/log"
)

// Security implements the 'apiKey' security type.
type Security struct {
	APISecurity *openapi3.SecurityScheme
	Log         log.ILogger
}

// Secure adds an example value from the API spec to the Authorization request header.
func (sec Security) Secure(req *http.Request) {
	switch sec.APISecurity.In {
	case "header":
		Header{Log: sec.Log, APISecurity: *sec.APISecurity}.Secure(req)

	case "query":
		Query{Log: sec.Log, APISecurity: *sec.APISecurity}.Secure(req)

	case "cookie":
		Cookie{Log: sec.Log, APISecurity: *sec.APISecurity}.Secure(req)

	default:
		fmt.Printf("\tUnknown API key location \"%s\"\n", sec.APISecurity.In)
	}
}
