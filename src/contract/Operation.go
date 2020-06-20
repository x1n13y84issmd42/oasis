package contract

import (
	"net/http"
)

// Operation ...
type Operation interface {
	ID() string
	Name() string
	Description() string
	Method() string
	Path() string
	CreateRequest() *http.Request
}

// OperationPrototype is a prototype implementation for operations.
type OperationPrototype struct {
	Operation
	Data DataProvider
}

// NewOperationPrototype create a new OperationPrototype instance.
func NewOperationPrototype() *OperationPrototype {
	return &OperationPrototype{
		Data: DataProvider{},
	}
}

// CreateRequest creates an http.Request instance and prepares it to make an API request.
func (op *OperationPrototype) CreateRequest() *http.Request {
	res, _ := http.NewRequest(op.Method(), "http://example.com", nil)
	op.Data.URL.Enrich(res)
	op.Data.Query.Enrich(res)
	op.Data.Headers.Enrich(res)
	return res
}
