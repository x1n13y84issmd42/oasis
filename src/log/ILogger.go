package log

import (
	"fmt"

	"github.com/x1n13y84issmd42/oasis/src/api"
	"github.com/xeipuuv/gojsonschema"
)

// ILogger - interface for test execution loggers.
type ILogger interface {
	NOMESSAGE(msg string, args ...interface{})

	Usage()
	Error(err error)
	LoadingSpec(path string)

	PrintOperations(ops []*api.Operation)
	TestingProject(p *api.ProjectInfo)
	TestingOperation(res *api.Operation)

	UsingHost(p *api.Host)
	UsingDefaultHost()
	HostNotFound(h string)

	UsingSecurity(sec api.ISecurity)
	SecurityHasNoData(sec api.ISecurity)

	Requesting(method string, url string)

	ParameterHasNoExample(paramName string, in string, container string)
	UsingParameterExample(paramName string, in string, container string)

	HeaderHasNoValue(header *api.Header)
	ResponseHasWrongStatus(schema *api.Response, actualStatus int)
	ResponseHasWrongContentType(schema *api.Response, actualCT string)

	OperationOK(res *api.Operation)
	OperationFail(res *api.Operation)
	OperationNotFound(op string)

	SchemaTesting(schema *api.Schema, data interface{})
	SchemaOK(schema *api.Schema)
	SchemaFail(schema *api.Schema, errors []gojsonschema.ResultError)

	XError(err error, style IStyle, tab TabFn)
}

// New creates a new logger based on the provided log style & level.
func New(style string, level int64) ILogger {
	switch style {
	case "plain":
		return NewPlain(level)

	case "festive":
		return NewFestive(level)
	}

	fmt.Printf("The \"%s\" log style is unknown.\nAvailable loggers are:\n", style)
	fmt.Println("\tplain - a simple text logger")
	fmt.Println("\tfestive - a nicer colorized logger")

	panic("No way.")
}
