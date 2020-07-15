package contract

import (
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

// Logger - interface for test execution loggers.
type Logger interface {
	NOMESSAGE(msg string, args ...interface{})
	Parameters(name string, params ParameterSource)

	Print(l int64, msg string, args ...interface{})

	Usage()
	Error(err error)
	LoadingSpec(path string)
	LoadingScript(path string)

	PrintOperations(ops OperationIterator)
	TestingProject(p ProjectInfo)
	TestingOperation(res Operation)

	UsingSecurity(sec Security)
	SecurityHasNoData(sec Security)

	Requesting(method string, url string)

	UsingParameterExample(paramName string, in string, container string, value string)

	Expecting(what string, v string)

	HeaderHasNoValue(hdr string)
	ResponseHasWrongStatus(expectedStatus int, actualStatus int)
	ResponseHasWrongContentType(expectedCT string, actualCT string)

	OperationOK()
	OperationFail()

	SchemaOK()
	SchemaFail(errors []gojsonschema.ResultError)

	ScriptExecutionStart(node string)

	XError(err error, style LogStyle, tab TabFn)
}

// LogStyle is an interface to log styling.
type LogStyle interface {
	Default(...interface{}) string
	URL(...interface{}) string
	Method(...interface{}) string
	Op(...interface{}) string
	OK(...interface{}) string
	Failure(...interface{}) string
	Success(...interface{}) string
	Error(...interface{}) string
	ID(...interface{}) string
	ValueExpected(...interface{}) string
	ValueActual(...interface{}) string
}

// TabFn produces indentation when printing nested errors.
type TabFn func(log Logger)

// Tab creates an indentation function.
func Tab(level uint) TabFn {
	return func(log Logger) {
		log.Print(1, ""+strings.Repeat("  ", int(level)))
	}
}

// Shift produces a new indentation function by adding a level to it.
func (fn TabFn) Shift() TabFn {
	return func(log Logger) {
		fn(log)
		log.Print(1, "  ")
	}
}
