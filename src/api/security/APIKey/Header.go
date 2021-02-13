package apikey

import (
	"net/http"

	"github.com/x1n13y84issmd42/oasis/src/contract"
)

// Header represents an API token security which expects a token as a header value.
type Header struct {
	Security
}

// Enrich adds an API key to the request's headers.
func (sec *Header) Enrich(req *http.Request, log contract.Logger) {
	log.UsingSecurity(sec)

	if v := sec.Value(); v != "" {
		req.Header.Add(sec.ParamName, v)
	} else {
		log.SecurityHasNoData(sec)
	}
}
