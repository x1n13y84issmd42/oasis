package security

import (
	"fmt"
	"net/http"

	"github.com/x1n13y84issmd42/goasis/src/api"
	"github.com/x1n13y84issmd42/goasis/src/log"
	APIKey "github.com/x1n13y84issmd42/goasis/src/test/security/APIKey"
	HTTP "github.com/x1n13y84issmd42/goasis/src/test/security/http"
)

// Security implements the OAS security mechanisms.
type Security struct {
	APISecurity *api.Security
	Log         log.ILogger
}

// NewSecurity creates a new Security instance.
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
		HTTP.Security{APISecurity: sec.APISecurity, Log: sec.Log}.Secure(req)

	case api.SecurityTypeAPIKey:
		APIKey.Security{APISecurity: sec.APISecurity, Log: sec.Log}.Secure(req)

	default:
		fmt.Printf("Unknown security type '%s'\n", sec.APISecurity.SecurityType)
	}
}
