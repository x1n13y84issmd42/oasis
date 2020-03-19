package apikey

import (
	"fmt"
	"net/http"

	"github.com/x1n13y84issmd42/oasis/src/api"
	"github.com/x1n13y84issmd42/oasis/src/log"
)

// Cookie represents an API token security which expects a token as a header value.
type Cookie struct {
	Log         log.ILogger
	APISecurity api.Security
}

// Secure adds an API key to the request's headers.
func (sec Cookie) Secure(req *http.Request) {
	if sec.APISecurity.Example == "" {
		fmt.Printf("The security \"%s\" contains no example to use in request.", sec.APISecurity)
	} else {
		req.AddCookie(&http.Cookie{Name: sec.APISecurity.ParamName, Value: sec.APISecurity.Example})
	}
}
