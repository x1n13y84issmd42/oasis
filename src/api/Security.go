package api

import (
	"net/http"
)

// ISecurity is an interface for security mechanisms.
type ISecurity interface {
	Secure(*http.Request)
	GetName() string
}
