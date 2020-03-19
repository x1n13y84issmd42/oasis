package apikey

import (
	"fmt"
	"net/http"

	"github.com/x1n13y84issmd42/oasis/src/api"
	"github.com/x1n13y84issmd42/oasis/src/log"
)

// Security implements the 'apiKey' security type.
type Security struct {
	APISecurity *api.Security
	Log         log.ILogger
}

// Secure adds an example value from the API spec to the Authorization request header.
func (sec Security) Secure(req *http.Request) {
	switch sec.APISecurity.In {
	case api.ParameterLocationHeader:
		Header{Log: sec.Log, APISecurity: *sec.APISecurity}.Secure(req)

	case api.ParameterLocationQuery:
		Query{Log: sec.Log, APISecurity: *sec.APISecurity}.Secure(req)

	case api.ParameterLocationCookie:
		Cookie{Log: sec.Log, APISecurity: *sec.APISecurity}.Secure(req)

	default:
		fmt.Printf("\tUnknown API key location \"%s\"\n", sec.APISecurity.In)
	}
}
