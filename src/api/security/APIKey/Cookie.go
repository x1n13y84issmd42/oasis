package apikey

import (
	"net/http"

	"github.com/x1n13y84issmd42/oasis/src/contract"
)

// Cookie represents an API token security which expects a token as a header value.
type Cookie struct {
	Security
}

// Enrich adds an API key to the request's headers.
func (sec Cookie) Enrich(req *http.Request, log contract.Logger) {
	log.UsingSecurity(sec)

	if sec.Value != "" {
		req.AddCookie(&http.Cookie{Name: sec.ParamName, Value: sec.Value})
	} else {
		log.SecurityHasNoData(sec)
	}
}
