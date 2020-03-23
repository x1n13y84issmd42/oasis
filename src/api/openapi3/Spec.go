package openapi3

import (
	"net/url"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/x1n13y84issmd42/oasis/src/api"
)

type Spec struct {
	OAS *openapi3.Swagger
}

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
	return &api.Operation{
		ID:          oasOp.OperationID,
		Name:        oasOp.Summary,
		Description: oasOp.Description,
		Request: api.Request{
			Method: method,
			URL:    spec.createURL(oasPath, oasPathItem, oasOp),
		},
	}
}

// createURL creates a fully qualified URL by joining
// the server host name with an operation path
// and replaceing path parameters with actual values from `example`.
func (spec *Spec) createURL(oasPath string, oasPathItem *openapi3.PathItem, oasOp *openapi3.Operation) url.URL {
	u, _ := url.Parse(oasPath)
	return *u
}

func (spec *Spec) GetProjectInfo() *api.ProjectInfo {
	return &api.ProjectInfo{
		Title:       spec.OAS.Info.Title,
		Description: spec.OAS.Info.Description,
		Version:     spec.OAS.Info.Version,
	}
}

func (spec *Spec) GetHost(name string) *api.Host {
	return nil
}

func (spec *Spec) GetDefaultHost() *api.Host {
	return nil
}
