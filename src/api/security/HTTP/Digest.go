package http

import (
	"net/http"
)

// Digest implements a Digest http authentication.
type Digest struct {
	Security
}

// Secure adds an example value from the API spec to the Authorization request header.
func (sec Digest) Secure(req *http.Request) {
	sec.Log.UsingSecurity(sec)

	if sec.Token != "" {
		req.Header["Authorization"] = append(req.Header["Authorization"], sec.Token)
	} else if sec.Username != "" {
		//TODO: implement client-side digest encoding
	} else {
		sec.Log.NOMESSAGE("The security \"%s\" contains no data to use in request.", sec.Name)
	}
}
