package openapi3

import (
	"sort"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/x1n13y84issmd42/oasis/src/api"
	"github.com/x1n13y84issmd42/oasis/src/contract"
	"github.com/x1n13y84issmd42/oasis/src/errors"
	"github.com/x1n13y84issmd42/oasis/src/params"
)

// Spec is an OAS3-backed API test spec.
type Spec struct {
	Log contract.Logger
	OAS *openapi3.Swagger
}

// Operations returns an iterable channel with operations.
func (spec *Spec) Operations() contract.OperationIterator {

	ch := make(contract.OperationIterator)

	addOp := func(oasOp *openapi3.Operation, method string, oasPath string, oasPathItem *openapi3.PathItem) {
		if oasOp != nil {
			ch <- spec.MakeOperation(method, oasOp, oasPath, oasPathItem)
		}
	}

	go func() {
		// First, sorting the paths lexicographically.
		paths := []string{}

		for oasPath := range spec.OAS.Paths {
			paths = append(paths, oasPath)
		}

		sort.Strings(paths)

		// Next, iterating over the sorted paths and feeding the ops to the channel.
		for _, oasPath := range paths {
			oasPathItem := spec.OAS.Paths[oasPath]

			addOp(oasPathItem.Get, "GET", oasPath, oasPathItem)
			addOp(oasPathItem.Post, "POST", oasPath, oasPathItem)
			addOp(oasPathItem.Put, "PUT", oasPath, oasPathItem)
			addOp(oasPathItem.Delete, "DELETE", oasPath, oasPathItem)
			addOp(oasPathItem.Patch, "PATCH", oasPath, oasPathItem)
			addOp(oasPathItem.Head, "HEAD", oasPath, oasPathItem)
			addOp(oasPathItem.Options, "OPTIONS", oasPath, oasPathItem)
			addOp(oasPathItem.Trace, "TRACE", oasPath, oasPathItem)
			addOp(oasPathItem.Connect, "CONNECT", oasPath, oasPathItem)
		}

		close(ch)
	}()

	return ch
}

// GetOperation returns a list of all available test operations from the spec.
func (spec *Spec) GetOperation(name string) contract.Operation {
	filterOp := func(oasOp *openapi3.Operation) bool {
		return (oasOp != nil && (oasOp.Summary == name || oasOp.OperationID == name))
	}

	for oasPath, oasPathItem := range spec.OAS.Paths {
		if filterOp(oasPathItem.Get) {
			return spec.MakeOperation("GET", oasPathItem.Get, oasPath, oasPathItem)
		}
		if filterOp(oasPathItem.Post) {
			return spec.MakeOperation("POST", oasPathItem.Post, oasPath, oasPathItem)
		}
		if filterOp(oasPathItem.Put) {
			return spec.MakeOperation("PUT", oasPathItem.Put, oasPath, oasPathItem)
		}
		if filterOp(oasPathItem.Delete) {
			return spec.MakeOperation("DELETE", oasPathItem.Delete, oasPath, oasPathItem)
		}
		if filterOp(oasPathItem.Patch) {
			return spec.MakeOperation("PATCH", oasPathItem.Patch, oasPath, oasPathItem)
		}
		if filterOp(oasPathItem.Head) {
			return spec.MakeOperation("HEAD", oasPathItem.Head, oasPath, oasPathItem)
		}
		if filterOp(oasPathItem.Options) {
			return spec.MakeOperation("OPTIONS", oasPathItem.Options, oasPath, oasPathItem)
		}
		if filterOp(oasPathItem.Connect) {
			return spec.MakeOperation("CONNECT", oasPathItem.Connect, oasPath, oasPathItem)
		}
		if filterOp(oasPathItem.Trace) {
			return spec.MakeOperation("TRACE", oasPathItem.Trace, oasPath, oasPathItem)
		}
	}

	return api.NoOperation(errors.NotFound("Operation", name, nil), spec.Log)
}

// MakeOperation creates an Operation instance from available spec data.
func (spec *Spec) MakeOperation(
	method string,
	oasOp *openapi3.Operation,
	oasPath string,
	oasPathItem *openapi3.PathItem,
) contract.Operation {
	op := &Operation{
		OperationPrototype: api.NewOperationPrototype(spec.Log),
		RequestMethod:      method,
		RequestPath:        oasPath,
		SpecOp:             oasOp,
		SpecPath:           oasPathItem,
	}

	op.Resolver = NewDataResolver(spec.Log, spec.OAS, op, &oasOp.Responses)
	op.OperationPrototype.Operation = op

	op.Data().URL = params.URL(oasPath, spec.Log)
	op.Data().URL.Load(PathParameterSource(&op.SpecPath.Parameters, "path"))
	op.Data().URL.Load(PathParameterSource(&op.SpecOp.Parameters, "op"))

	op.Data().Query = params.Query(spec.Log)
	op.Data().Query.Load(QueryParameterSource(&op.SpecPath.Parameters, "path"))
	op.Data().Query.Load(QueryParameterSource(&op.SpecOp.Parameters, "op"))

	op.Data().Headers = params.Headers(spec.Log)
	op.Data().Headers.Load(HeadersParameterSource(&op.SpecPath.Parameters, "path"))
	op.Data().Headers.Load(HeadersParameterSource(&op.SpecOp.Parameters, "op"))

	op.Data().Body = params.Body(spec.Log)

	requireParameters := func(p *openapi3.Parameter) {
		switch p.In {
		case "path":
			op.Data().URL.Require(p.Name)
			break
		case "query":
			op.Data().Query.Require(p.Name)
			break
		case "headers":
			op.Data().Headers.Require(p.Name)
			break
		}
	}

	spec.IterateOverRequiredParameters(&op.SpecPath.Parameters, requireParameters)
	spec.IterateOverRequiredParameters(&op.SpecOp.Parameters, requireParameters)

	return op
}

// IterateOverRequiredParameters iterates over items in the provided parameter list
// and invokes the handler function for every required one.
func (spec *Spec) IterateOverRequiredParameters(params *openapi3.Parameters, handler func(*openapi3.Parameter)) {
	for _, specP := range *params {
		if specP == nil || specP.Value == nil || !specP.Value.Required {
			continue
		}

		handler(specP.Value)
	}
}

// Title return project title.
func (spec *Spec) Title() string {
	return spec.OAS.Info.Title
}

// Description return project description.
func (spec *Spec) Description() string {
	return spec.OAS.Info.Description
}

// Version return project version.
func (spec *Spec) Version() string {
	return spec.OAS.Info.Version
}
