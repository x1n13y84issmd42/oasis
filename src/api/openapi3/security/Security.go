package security

import (
	"fmt"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	APIKey "github.com/x1n13y84issmd42/oasis/src/api/openapi3/security/APIKey"
	HTTP "github.com/x1n13y84issmd42/oasis/src/api/openapi3/security/http"
	"github.com/x1n13y84issmd42/oasis/src/log"
)

// Security implements the OAS security mechanisms.
type Security struct {
	APISecurity *openapi3.SecurityScheme
	Log         log.ILogger
}

// NewSecurity creates a new Security instance.
func NewSecurity(apiSec *openapi3.SecurityScheme, log log.ILogger) Security {
	return Security{
		APISecurity: apiSec,
		Log:         log,
	}
}

// Secure modifies the passed Request instance so it contains
// authentication credentials needed to perform an operation.
func (sec Security) Secure(req *http.Request) {
	switch sec.APISecurity.Type {
	case "http":
		HTTP.Security{APISecurity: sec.APISecurity, Log: sec.Log}.Secure(req)

	case "apiKey":
		APIKey.Security{APISecurity: sec.APISecurity, Log: sec.Log}.Secure(req)

	default:
		fmt.Printf("Unknown security type '%s'\n", sec.APISecurity.Type)
	}
}
