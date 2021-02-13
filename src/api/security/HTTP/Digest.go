package http

import (
	"net/http"

	"github.com/x1n13y84issmd42/oasis/src/contract"
)

// Digest implements a Digest HTTP authentication.
type Digest struct {
	Security
}

// Enrich adds an example value from the API spec to the Authorization request header.
func (sec *Digest) Enrich(req *http.Request, log contract.Logger) {
	log.UsingSecurity(sec)

	if t := sec.Token(); t != "" {
		req.Header["Authorization"] = append(req.Header["Authorization"], t)
	} else if u := sec.Username(); u != "" {
		//TODO: implement client-side digest encoding
	} else {
		log.SecurityHasNoData(sec)
	}
}
