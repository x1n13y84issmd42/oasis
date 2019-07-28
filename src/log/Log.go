package log

import (
	"github.com/x1n13y84issmd42/goasis/src/api"
	"github.com/x1n13y84issmd42/goasis/src/utility"
)

// ILogger - interface for execution loggers
type ILogger interface {
	Error(err error)
	TestingProject(p *api.ProjectInfo)
	UsingHost(p *api.Host)
	UsingDefaultHost()
	TestingOperation(res *api.Operation)
	UsingSecurity(sec *api.Security)
	UsingRequest(req *api.Request)
	UsingResponse(req *api.Response)

	Overriding(what string)

	PropertyHasNoValue(prop *api.Property, ctx *utility.Context)
	PropertyHasWrongType(prop *api.Property, ctx *utility.Context)
	HeaderHasNoValue(schema *api.Header)
	HeaderHasWrongType(schema *api.Header)
	ResponseHasWrongStatus(schema *api.Response, actualStatus int)
	ResponseHasWrongContentType(schema *api.Response, actualCT string)
	ResponseExpectedArray(schema *api.Response)
	ResponseExpectedObject(schema *api.Response)

	OperationOK(res *api.Operation)
	OperationFail(res *api.Operation)
	SchemaOK(schema *api.Schema)
	SchemaFail(schema *api.Schema)
}
