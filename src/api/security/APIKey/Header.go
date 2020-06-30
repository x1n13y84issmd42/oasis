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
func (sec Header) Enrich(req *http.Request, log contract.Logger) {
	log.UsingSecurity(sec)

	if sec.Value != "" {
		req.Header.Add(sec.ParamName, sec.Value)
	} else {
		log.SecurityHasNoData(sec)
	}
}
