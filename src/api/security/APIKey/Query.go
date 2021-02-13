package apikey

import (
	"net/http"

	"github.com/x1n13y84issmd42/oasis/src/contract"
)

// Query represents an API token security which expects a token as a query parameter.
type Query struct {
	Security
}

// Enrich adds an API key to the request's query.
func (sec *Query) Enrich(req *http.Request, log contract.Logger) {
	log.UsingSecurity(sec)

	if v := sec.Value(); v != "" {
		q := req.URL.Query()
		q.Add(sec.ParamName, v)
		req.URL.RawQuery = q.Encode()
	} else {
		log.SecurityHasNoData(sec)
	}
}
