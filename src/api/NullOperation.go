package api

import (
	"net/http"

	"github.com/x1n13y84issmd42/oasis/src/contract"
)

// NullOperation is used whenever we can't have a real one.
// Reports the contained error on every method call.
type NullOperation struct {
	Log   contract.Logger
	Error error
}

// CreateRequest ...
func (op NullOperation) CreateRequest() *http.Request {
	return nil
}

// ID ...
func (op NullOperation) ID() string {
	return ""
}

// Name ...
func (op NullOperation) Name() string {
	return ""
}

// Description ...
func (op NullOperation) Description() string {
	return ""
}

// Method ...
func (op NullOperation) Method() string {
	return ""
}

// Path ...
func (op NullOperation) Path() string {
	return ""
}
