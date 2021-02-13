package http

import (
	"errors"
	"net/http"
	"strings"

	"github.com/x1n13y84issmd42/oasis/src/api"
	"github.com/x1n13y84issmd42/oasis/src/contract"
	"github.com/x1n13y84issmd42/oasis/src/params"
)

// Security implements the 'http' security type.
type Security struct {
	Name     string
	Token    contract.ParameterAccess
	Username contract.ParameterAccess
	Password contract.ParameterAccess
	Log      contract.Logger
}

// WWWAuthenticateQoP is a type for the Quality of Protection value.
// Used in the HTTP Digest auth.
type WWWAuthenticateQoP string

// QoP types.
const (
	WWWAuthenticateQoPAuth     = WWWAuthenticateQoP("auth")
	WWWAuthenticateQoPAuthIntl = WWWAuthenticateQoP("auth-intl")
)

// WWWAuthenticate is a representation of the Www-Authenticate HTTP response header.
type WWWAuthenticate struct {
	Realm  string
	Nonce  string
	CNonce string
	QoP    WWWAuthenticateQoP
}

// New creates a new HTTP security.
func New(name string, scheme string, token string, username string, password string, logger contract.Logger) contract.Security {
	switch scheme {
	case "basic":
		return &Basic{
			Security{
				Name:     name,
				Token:    params.Value(token),
				Log:      logger,
				Username: params.Value(username),
				Password: params.Value(password),
			},
		}

	case "digest":
		return &Digest{
			Security{
				Name:  name,
				Token: params.Value(token),
				Log:   logger,
			},
		}
	}

	return api.NoSecurity(errors.New("Unknown security scheme '"+scheme+"'"), logger)
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
// Digest realm="Oasis",nonce="61b6948856629ad7fd3da9d6179393ec",qop="auth",opaque="f9a0f11abf3f6710d22c5a2aa65e19036"
// The function returns these 'realm', 'nonce' and other directives as a map.
//TODO: it fails when quoted values have commas, like "qop"
// Digest realm="Oasis",nonce="61b6948856629ad7fd3da9d6179393ec",qop="auth,auth-int",opaque="f9a0f11abf3f6710d22c5a2aa65e19036"
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

// GetName returns name.
func (sec Security) GetName() string {
	return sec.Name
}

// SetValue does nothing.
func (sec *Security) SetValue(v contract.ParameterAccess) {
}

// SetToken sets Token.
func (sec *Security) SetToken(v contract.ParameterAccess) {
	sec.Token = v
}

// SetUsername sets Username.
func (sec *Security) SetUsername(v contract.ParameterAccess) {
	sec.Username = v
}

// SetPassword sets Password.
func (sec *Security) SetPassword(v contract.ParameterAccess) {
	sec.Password = v
}
