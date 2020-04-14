package apikey

import (
	"net/http"
)

// Header represents an API token security which expects a token as a header value.
type Header struct {
	Security
}

// Secure adds an API key to the request's headers.
func (sec Header) Secure(req *http.Request) {
	sec.Log.UsingSecurity(sec)

	if sec.Value != "" {
		req.Header.Add(sec.ParamName, sec.Value)
	} else {
		sec.Log.SecurityHasNoData(sec)
	}
}
