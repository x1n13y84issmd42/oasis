package apikey

import (
	"fmt"
	"net/http"

	"github.com/x1n13y84issmd42/goasis/src/api"
	"github.com/x1n13y84issmd42/goasis/src/log"
)

// Query represents an API token security which expects a token as a query parameter.
type Query struct {
	Log         log.ILogger
	APISecurity api.Security
}

// Secure adds an API key to the request's query.
func (sec Query) Secure(req *http.Request) {
	if sec.APISecurity.Example == "" {
		fmt.Printf("The security \"%s\" contains no example to use in request.\n", sec.APISecurity)
	} else {
		q := req.URL.Query()
		q.Add(sec.APISecurity.ParamName, sec.APISecurity.Example)
		req.URL.RawQuery = q.Encode()
	}
}
