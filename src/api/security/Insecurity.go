package security

import (
	"net/http"

	"github.com/x1n13y84issmd42/oasis/src/contract"
)

// Empty doesn't implement any security.
type Empty struct {
	Log contract.Logger
}

// GetName returns nothing.
func (sec Empty) GetName() string {
	return "Insecurity"
}

// SetValue does nothing for Insecurity.
func (sec *Empty) SetValue(v contract.ParameterAccess) {
}

// SetToken does nothing for Insecurity.
func (sec *Empty) SetToken(v contract.ParameterAccess) {
}

// SetUsername does nothing for Insecurity.
func (sec *Empty) SetUsername(v contract.ParameterAccess) {
}

// SetPassword does nothing for Insecurity.
func (sec *Empty) SetPassword(v contract.ParameterAccess) {
}

// Enrich does nothing for Insecurity.
func (sec *Empty) Enrich(req *http.Request, log contract.Logger) {
	sec.Log.UsingSecurity(sec)
}

// Insecurity creates a new Empty security instance.
func Insecurity(log contract.Logger) contract.Security {
	return &Empty{
		Log: log,
	}
}
