package log

import "github.com/x1n13y84issmd42/oasis/src/api"

// ILogger - interface for execution loggers
type ILogger interface {
	Usage()
	PrintOperations(ops []*api.Operation)
	Error(err error)
	TestingProject(p *api.ProjectInfo)
	TestingOperation(res *api.Operation)

	UsingHost(p *api.Host)
	UsingDefaultHost()
	HostNotFound(h string)

	// UsingSecurity(sec *api.Security)
	// UsingRequest(req *api.Request)
	UsingResponse(req *api.Response)

	Overriding(what string)
	Requesting(url string)

	// ParameterHasNoExample(param *api.Parameter, container string)
	// UsingParameterExample(param *api.Parameter, container string)

	// PropertyHasNoValue(prop *api.Property, ctx *utility.Context)
	// PropertyHasWrongType(prop *api.Property, ctx *utility.Context)
	// HeaderHasNoValue(schema *api.Header)
	// HeaderHasWrongType(schema *api.Header)
	ResponseHasWrongStatus(schema *api.Response, actualStatus int)
	ResponseHasWrongContentType(schema *api.Response, actualCT string)
	ResponseExpectedArray(schema *api.Response)
	ResponseExpectedObject(schema *api.Response)
	ResponseNotFound(CT string, status int)

	OperationOK(res *api.Operation)
	OperationFail(res *api.Operation)
	OperationNotFound(op string)

	// SchemaOK(schema *api.Schema)
	// SchemaFail(schema *api.Schema)
	// UnknownSchemaDataType(schema *api.Schema)
	// SchemaExpectedBoolean(schema *api.Schema, v interface{})
	// SchemaExpectedNumber(schema *api.Schema, v interface{})
	// SchemaExpectedInteger(schema *api.Schema, v interface{})
	// SchemaExpectedString(schema *api.Schema, v interface{})
	// SchemaExpectedArray(schema *api.Schema, v interface{})
	// SchemaExpectedObject(schema *api.Schema, v interface{})
}
