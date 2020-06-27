package api

import (
	"net/http"

	"github.com/x1n13y84issmd42/oasis/src/contract"
)

// OperationPrototype is a prototype implementation for operations.
type OperationPrototype struct {
	contract.EntityTrait
	contract.Operation

	data contract.OperationData
}

// NewOperationPrototype create a new OperationPrototype instance.
func NewOperationPrototype(log contract.Logger) *OperationPrototype {
	return &OperationPrototype{
		EntityTrait: contract.Entity(log),
		data:        contract.OperationData{},
	}
}

// GetRequest creates an http.Request instance and prepares it to make an API request.
func (op *OperationPrototype) GetRequest() *http.Request {
	res, _ := http.NewRequest(op.Method(), op.data.URL.String(), nil)

	//TODO: these should come from the spec along with the op and be treated
	// as any other request enrichment thing (like security).
	// op.Data().Query.Enrich(res)
	// op.Data().Headers.Enrich(res)
	return res
}

// Data creates an http.Request instance and prepares it to make an API request.
func (op *OperationPrototype) Data() *contract.OperationData {
	return &op.data
}
