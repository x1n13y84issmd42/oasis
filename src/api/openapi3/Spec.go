package openapi3

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/x1n13y84issmd42/oasis/src/api"
)

// Spec is an OAS3-backed API test spec.
type Spec struct {
	OAS        *openapi3.Swagger
	Parameters *api.OperationParameters
}

// GetOperations returns a list of all available test operations from the spec.
func (spec *Spec) GetOperations() []*api.Operation {
	ops := []*api.Operation{}

	addOp := func(oasOp *openapi3.Operation, method string, oasPath string, oasPathItem *openapi3.PathItem) {
		if oasOp != nil {
			ops = append(ops, spec.makeOperation(method, oasOp, oasPath, oasPathItem))
		}
	}

	for oasPath, oasPathItem := range spec.OAS.Paths {
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

	return ops
}

func (spec *Spec) makeOperation(method string, oasOp *openapi3.Operation, oasPath string, oasPathItem *openapi3.PathItem) *api.Operation {
	specPath := spec.CreatePath(oasPath, oasPathItem, oasOp)
	specRequests := []*api.Request{}

	if oasOp.RequestBody != nil && oasOp.RequestBody.Value != nil {
		for oasCT, oasMT := range oasOp.RequestBody.Value.Content {
			specRequests = append(specRequests, spec.MakeRequest(method, specPath, oasOp, oasPathItem, oasCT, oasMT))
		}
	} else {
		specRequests = append(specRequests, spec.MakeRequest(method, specPath, oasOp, oasPathItem, "", nil))
	}

	return &api.Operation{
		ID:          oasOp.OperationID,
		Name:        oasOp.Summary,
		Description: oasOp.Description,
		Method:      method,
		Path:        specPath,
		Requests:    specRequests,
	}
}

// MakeRequest creates an api.Rrequest instance from available operation spec data.
func (spec *Spec) MakeRequest(
	method string,
	specPath string,
	oasOp *openapi3.Operation,
	oasPathItem *openapi3.PathItem,
	oasCT string,
	oasMT *openapi3.MediaType,
) *api.Request {

	specReq := &api.Request{}
	specReq.Method = method
	specReq.Path = specPath

	return specReq
}

// CreatePath creates an operation path with parameters replaced by actual values from `example`.
// Examples from operation-level parameters override examples from the path-level ones.
func (spec *Spec) CreatePath(oasPath string, oasPathItem *openapi3.PathItem, oasOp *openapi3.Operation) string {
	return oasPath
}

// GetProjectInfo returns project info from the spec.info object.
func (spec *Spec) GetProjectInfo() *api.ProjectInfo {
	return &api.ProjectInfo{
		Title:       spec.OAS.Info.Title,
		Description: spec.OAS.Info.Description,
		Version:     spec.OAS.Info.Version,
	}
}

// GetHost returns an API host by requested description
// from the spec.servers list.
func (spec *Spec) GetHost(name string) *api.Host {
	for _, oasServer := range spec.OAS.Servers {
		if oasServer.Description == name {
			return &api.Host{
				Name:        "Default",
				Description: oasServer.Description,
				URL:         oasServer.URL,
			}
		}
	}
	return nil
}

// GetDefaultHost returns the fisr host from the spec.servers list as default.
func (spec *Spec) GetDefaultHost() *api.Host {
	if len(spec.OAS.Servers) > 0 {
		return &api.Host{
			Name:        "Default",
			Description: spec.OAS.Servers[0].Description,
			URL:         spec.OAS.Servers[0].URL,
		}
	}
	return nil
}
