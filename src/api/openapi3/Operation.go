package openapi3

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/x1n13y84issmd42/oasis/src/api"
	"github.com/x1n13y84issmd42/oasis/src/contract"
)

// Operation provides access to OAS3-specific API data.
type Operation struct {
	*api.OperationPrototype

	RequestMethod string
	RequestPath   string
	SpecPath      *openapi3.PathItem
	SpecOp        *openapi3.Operation

	Resolver *DataResolver
}

// ID returns an operation ID.
func (op *Operation) ID() string {
	return op.SpecOp.OperationID
}

// Name returns an operation name.
func (op *Operation) Name() string {
	return op.SpecOp.Summary
}

// Description returns an operation description.
func (op *Operation) Description() string {
	return op.SpecOp.Description
}

// Method returns an operation request method.
func (op *Operation) Method() string {
	return op.RequestMethod
}

// Path returns an operation request path.
func (op *Operation) Path() string {
	return op.RequestPath
}

// Resolve returns a DataResolver instance which is used
// to resolve data based on user input.
func (op *Operation) Resolve() contract.DataResolver {
	return op.Resolver
}
