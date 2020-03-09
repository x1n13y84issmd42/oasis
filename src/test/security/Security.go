package security

import (
	"fmt"
	"net/http"

	"github.com/x1n13y84issmd42/goasis/src/api"
	"github.com/x1n13y84issmd42/goasis/src/log"
	"github.com/x1n13y84issmd42/goasis/src/test/security/HTTP"
)

// Security implements the OAS security mechanisms.
type Security struct {
	APISecurity *api.Security
	Log         log.ILogger
}

func NewSecurity(apiSec *api.Security, log log.ILogger) Security {
	return Security{
		APISecurity: apiSec,
		Log:         log,
	}
}

// Secure modifies the passed Request instance so it contains
// authentication credentials needed to perform an operation.
func (sec Security) Secure(req *http.Request) {
	switch sec.APISecurity.SecurityType {
	case api.SecurityTypeHTTP:
		HTTP.HTTPSecurity{sec.APISecurity, sec.Log}.Secure(req)

	default:
		fmt.Printf("Unknown security type '%s'\n", sec.APISecurity.SecurityType)
	}
}
