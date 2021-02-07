package api

import (
	"net/http"

	"github.com/x1n13y84issmd42/oasis/src/contract"
)

// OperationPrototype is a prototype implementation for operations.
type OperationPrototype struct {
	contract.EntityTrait
	contract.Operation

	data   contract.OperationData
	result contract.OperationResult
}

// NewOperationPrototype create a new OperationPrototype instance.
func NewOperationPrototype(log contract.Logger) *OperationPrototype {
	return &OperationPrototype{
		EntityTrait: contract.Entity(log.Clone()),
		data:        contract.OperationData{},
		result: contract.OperationResult{
			Success: true,
		},
	}
}

// GetRequest creates an http.Request instance and prepares it to make an API request.
func (op *OperationPrototype) GetRequest() (*http.Request, error) {
	return http.NewRequest(op.Method(), op.data.URL.String(), nil)
}

// Data returns a pointer to the internal data object.
func (op *OperationPrototype) Data() *contract.OperationData {
	return &op.data
}

// Result returns a pointer to the internal result object.
func (op *OperationPrototype) Result() *contract.OperationResult {
	return &op.result
}

// GetLogger returns logger.
func (op *OperationPrototype) GetLogger() contract.Logger {
	return op.EntityTrait.GetLogger()
}
