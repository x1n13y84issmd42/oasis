package openapi3

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/x1n13y84issmd42/oasis/src/api"
	"github.com/x1n13y84issmd42/oasis/src/contract"
	"github.com/x1n13y84issmd42/oasis/src/test/expect"
)

// Operation provides access to OAS3-specific API data.
type Operation struct {
	*api.OperationPrototype

	RequestMethod string
	RequestPath   string
	SpecPath      *openapi3.PathItem
	SpecOp        *openapi3.Operation

	Resolver DataResolver
}

// ID ...
func (op *Operation) ID() string {
	return op.SpecOp.OperationID
}

// Name ...
func (op *Operation) Name() string {
	return op.SpecOp.Summary
}

// Description ...
func (op *Operation) Description() string {
	return op.SpecOp.Description
}

// Method ...
func (op *Operation) Method() string {
	return op.RequestMethod
}

// Path ...
func (op *Operation) Path() string {
	return op.RequestPath
}

// GetValidator returns a validator instance used to test operation result.
func (op *Operation) GetValidator() contract.Validator {
	v := api.Validator{}
	v.Expect(expect.Status(400))
	return &v
}

// Resolve returns a DataResolver instance which is used
// to resolve data based on user input.
func (op *Operation) Resolve() contract.DataResolver {
	return &op.Resolver
}
