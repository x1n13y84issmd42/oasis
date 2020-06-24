package api

import (
	"net/http"

	"github.com/x1n13y84issmd42/oasis/src/contract"
)

// NullOperation is used whenever we can't have a real one.
// Reports the contained error on every method call.
type NullOperation struct {
	contract.NullObjectPrototype
}

// NoOperation creates a new NullOperation instance.
func NoOperation(err error, log contract.Logger) *NullOperation {
	return &NullOperation{
		NullObjectPrototype: contract.NullObject(err, log),
	}
}

// CreateRequest ...
func (op NullOperation) CreateRequest() *http.Request {
	op.Report()
	return nil
}

// ID ...
func (op NullOperation) ID() string {
	op.Report()
	return ""
}

// Name ...
func (op NullOperation) Name() string {
	op.Report()
	return ""
}

// Description ...
func (op NullOperation) Description() string {
	op.Report()
	return ""
}

// Method ...
func (op NullOperation) Method() string {
	op.Report()
	return ""
}

// Path ...
func (op NullOperation) Path() string {
	op.Report()
	return ""
}

// Data ...
func (op NullOperation) Data() *contract.OperationData {
	op.Report()
	return nil
}

// Host ...
func (op NullOperation) Host(h string) contract.ParameterSource {
	op.Report()
	return nil
}
