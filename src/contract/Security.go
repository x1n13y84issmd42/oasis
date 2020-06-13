package contract

import (
	"net/http"
)

// Security is an interface for security mechanisms.
type Security interface {
	Secure(*http.Request)
	GetName() string
}
