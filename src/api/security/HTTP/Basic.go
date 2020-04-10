package http

import (
	"net/http"
)

// Basic implements a Basic HTTP authentication.
type Basic struct {
	Security
}

// Secure adds an example value from the API spec to the Authorization request header.
func (sec Basic) Secure(req *http.Request) {
	sec.Log.UsingSecurity(sec)

	if sec.Token != "" {
		req.Header["Authorization"] = append(req.Header["Authorization"], sec.Token)
	} else if sec.Username != "" {
		//TODO: implement client-side basic encoding
	} else {
		sec.Log.NOMESSAGE("The security \"%s\" contains no data to use in request.", sec.Name)
	}
}
