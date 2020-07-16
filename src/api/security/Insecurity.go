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

// Enrich does nothing.
func (sec Empty) Enrich(req *http.Request, log contract.Logger) {
	sec.Log.UsingSecurity(sec)
}

// Insecurity creates a new Empty security instance.
func Insecurity(log contract.Logger) contract.Security {
	return &Empty{
		Log: log,
	}
}
