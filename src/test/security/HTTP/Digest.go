package HTTP

import (
	"net/http"

	"github.com/x1n13y84issmd42/goasis/src/api"
	"github.com/x1n13y84issmd42/goasis/src/log"
)

// Digest implements a Digest HTTP authentication.
type Digest struct {
	APISecurity *api.Security
	Log         log.ILogger
	Auth        WWWAuthenticate
}

// Secure adds an example value from the API spec to the Authorization request header.
func (sec Digest) Secure(req *http.Request) {
	req.Header["Authorization"] = append(req.Header["Authorization"], sec.APISecurity.Example)
}
