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

// GetRequest reports an error.
func (op NullOperation) GetRequest() (*http.Request, error) {
	op.Report()
	return nil, nil
}

// ID reports an error.
func (op NullOperation) ID() string {
	op.Report()
	return ""
}

// Name reports an error.
func (op NullOperation) Name() string {
	op.Report()
	return ""
}

// Description reports an error.
func (op NullOperation) Description() string {
	op.Report()
	return ""
}

// Method reports an error.
func (op NullOperation) Method() string {
	op.Report()
	return ""
}

// Path reports an error.
func (op NullOperation) Path() string {
	op.Report()
	return ""
}

// Data reports an error.
func (op NullOperation) Data() *contract.OperationData {
	op.Report()
	return nil
}

// Resolve reports an error.
func (op NullOperation) Resolve() contract.DataResolver {
	op.Report()
	return nil
}

// Result reports an error.
func (op NullOperation) Result() *contract.OperationResult {
	op.Report()
	return nil
}
