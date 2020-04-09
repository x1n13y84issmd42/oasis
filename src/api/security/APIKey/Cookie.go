package apikey

import (
	"fmt"
	"net/http"
)

// Cookie represents an API token security which expects a token as a header value.
type Cookie struct {
	Security
}

// Secure adds an API key to the request's headers.
func (sec Cookie) Secure(req *http.Request) {
	sec.Log.UsingSecurity(sec)

	if sec.Value != "" {
		req.AddCookie(&http.Cookie{Name: sec.ParamName, Value: sec.Value})
	} else {
		fmt.Printf("The security \"%s\" contains no data to use in request.", sec.Name)
	}
}
