package http

import (
	"encoding/base64"
	"net/http"

	"github.com/x1n13y84issmd42/oasis/src/contract"
)

// Basic implements a Basic HTTP authentication.
type Basic struct {
	Security
}

// Enrich adds an example value from the API spec to the Authorization request header.
func (sec *Basic) Enrich(req *http.Request, log contract.Logger) {
	log.UsingSecurity(sec)

	if t := sec.Token(); t != "" {
		req.Header["Authorization"] = append(req.Header["Authorization"], t)
	} else if u := sec.Username(); u != "" {
		token := base64.StdEncoding.EncodeToString([]byte(u + ":" + sec.Password()))
		token = "Basic " + token
		req.Header["Authorization"] = append(req.Header["Authorization"], token)
	} else {
		log.SecurityHasNoData(sec)
	}
}
