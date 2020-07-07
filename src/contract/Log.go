package contract

import (
	"strings"
)

// Logger - interface for test execution loggers.
type Logger interface {
	NOMESSAGE(msg string, args ...interface{})

	Print(l int64, msg string, args ...interface{})

	Usage()
	Error(err error)
	LoadingSpec(path string)

	PrintOperations(ops OperationIterator)
	TestingProject(p ProjectInfo)
	TestingOperation(res Operation)

	UsingSecurity(sec Security)
	SecurityHasNoData(sec Security)

	Requesting(method string, url string)

	// ParameterHasNoExample(paramName string, in string, container string)
	UsingParameterExample(paramName string, in string, container string, value string)

	HeaderHasNoValue(hdr string)
	ResponseHasWrongStatus(expectedStatus int, actualStatus int)
	ResponseHasWrongContentType(expectedCT string, actualCT string)

	OperationOK()
	OperationFail()
	OperationNotFound(op string)

	// SchemaTesting(schema *api.Schema, data interface{})
	// SchemaOK(schema *api.Schema)
	// SchemaFail(schema *api.Schema, errors []gojsonschema.ResultError)

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
