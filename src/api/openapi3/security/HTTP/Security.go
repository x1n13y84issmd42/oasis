package http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/x1n13y84issmd42/oasis/src/log"
)

// Security implements the 'http' security type.
type Security struct {
	APISecurity *openapi3.SecurityScheme
	Log         log.ILogger
}

// WWWAuthenticateQoP is a type for the Quality of Protection value.
// Used in the HTTP Digest auth.
type WWWAuthenticateQoP string

// QoP types.
const (
	WWWAuthenticateQoPAuth     = "auth"
	WWWAuthenticateQoPAuthIntl = "auth-intl"
)

// WWWAuthenticate is a representation of the Www-Authenticate HTTP response header.
type WWWAuthenticate struct {
	Realm  string
	Nonce  string
	CNonce string
	QoP    WWWAuthenticateQoP
}

// Secure adds an example value from the API spec to the Authorization request header.
func (sec Security) Secure(req *http.Request) {
	auth := sec.Probe(req)

	switch sec.APISecurity.Scheme {
	case "basic":
		Basic{sec.APISecurity, sec.Log, auth}.Secure(req)

	case "digest":
		Digest{sec.APISecurity, sec.Log, auth}.Secure(req)

	default:
		fmt.Printf("Unknown security scheme '%s'\n", sec.APISecurity.Scheme)
	}
}

// Probe makes a request to a URL which is (supposedly) protected
// by an HTTP Basic or Digest authentication scheme in order to obtain an authentication
// request from the server.
func (sec Security) Probe(req *http.Request) (auth WWWAuthenticate) {
	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	probeReq, _ := http.NewRequest(req.Method, req.URL.String(), nil)

	resp, _ := client.Do(probeReq)

	if resp != nil {
		if resp.Header["Www-Authenticate"] != nil && len(resp.Header["Www-Authenticate"]) > 0 {
			authMap := sec.ParseWWWAuthenticate(resp.Header["Www-Authenticate"][0])

			auth = WWWAuthenticate{
				Realm:  authMap["realm"],
				Nonce:  authMap["nonce"],
				CNonce: authMap["cnonce"],
				QoP:    WWWAuthenticateQoP(authMap["qop"]),
			}
		}
	} else {
		// fmt.Printf("No response from the %s URL\n", req.URL.String())
	}

	return
}

// ParseWWWAuthenticate parses the Www-Authenticate header value.
// A typical header looks something like this:
// Digest realm="Oasis",nonce="61b6948856629ad7fd3da9d6179393ec",qop="auth,auth-int",opaque="f9a0f11abf3f6710d22c5a2aa65e19036"
// The function returns these 'realm', 'nonce' and other directives as a map.
func (sec Security) ParseWWWAuthenticate(header string) map[string]string {
	directives := strings.Split(header, ",")

	res := make(map[string]string)

	for _, directive := range directives {
		directivePair := strings.Split(directive, "=")

		if directivePair[1][:1] == "\"" {
			directivePair[1] = directivePair[1][1:]
		}

		if directivePair[1][len(directivePair[1])-1:] == "\"" {
			directivePair[1] = directivePair[1][:len(directivePair[1])-1]
		}

		res[directivePair[0]] = directivePair[1]
	}

	return res
}
