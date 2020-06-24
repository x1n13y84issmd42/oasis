package openapi3

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/x1n13y84issmd42/oasis/src/api"
	"github.com/x1n13y84issmd42/oasis/src/contract"
	"github.com/x1n13y84issmd42/oasis/src/errors"
	"github.com/x1n13y84issmd42/oasis/src/params"
)

// Operation provides access to OAS3-specific API data.
type Operation struct {
	*api.OperationPrototype

	RequestMethod string
	RequestPath   string
	SpecPath      *openapi3.PathItem
	SpecOp        *openapi3.Operation
	SpecHosts     *openapi3.Servers
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

// Host attempts to come up with a valid host name based on the input
// host hint, which is a spec-specific host identifier.
// F.e. it's a name for OAS3, but may be index or a literal host for other
// spec standards.
func (op *Operation) Host(hostHint string) contract.ParameterSource {
	var host *string
	if hostHint == "" {
		if len(*op.SpecHosts) > 0 {
			host = &(*op.SpecHosts)[0].URL
		}
	} else {
		for _, oasServer := range *op.SpecHosts {
			if oasServer.Description == hostHint {
				host = &oasServer.URL
			}
		}
	}

	if host != nil {
		src := params.NewMemoryParameterSource()
		src.Add(params.KeyHost, *host)
		return src
	}

	return params.NoSource(errors.NotFound("Host", hostHint, nil), op.Log)

}
