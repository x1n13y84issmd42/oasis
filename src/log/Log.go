package log

import (
	"fmt"
	"strings"

	"github.com/x1n13y84issmd42/oasis/src/api"
	"github.com/x1n13y84issmd42/oasis/src/errors"
	"github.com/xeipuuv/gojsonschema"
)

// Log is a base type for loggers.
type Log struct {
	Level int64
	Style IStyle
}

// Print prints.
func (log Log) Print(l int64, msg string, args ...interface{}) {
	if l <= log.Level {
		fmt.Printf(msg, args...)
	}
}

// Println prints and adds a newline.
func (log Log) Println(l int64, msg string, args ...interface{}) {
	log.Print(l, msg+"\n", args...)
}

// NOMESSAGE is a default and temporary print function to use when you don't have
// a dedidated message function in the logger. It is meant to be replaced eventually
// with a proper logging method or removed. Henec the indicating name.
func (log Log) NOMESSAGE(msg string, args ...interface{}) {
	log.Println(1, "\t"+msg, args...)
}

// TabFn produces indentation when printing nested errors.
type TabFn func(log Log)

// Tab creates an indentation function.
func Tab(level uint) TabFn {
	return func(log Log) {
		log.Print(1, ""+strings.Repeat("  ", int(level)))
	}
}

// Shift produces a new indentation function by adding a level to it.
func (fn TabFn) Shift() TabFn {
	return func(log Log) {
		fn(log)
		log.Print(1, "  ")
	}
}

// Error outputs errors. It accepts both built-in errors as well as errors.IError instances.
// The latter carry their cause error, if any, and those will be recursively printed.
func (log Log) Error(err error) {
	log.XError(err, log.Style, Tab(0))
	log.Println(1, "")
}

// XError is an internal error handling function. It handles both
// built-in errors and errors.IError instances.
func (log Log) XError(err error, style IStyle, tab TabFn) {
	tab(log)
	log.Println(1, "%s", style.Error(err.Error()))

	if xerr, ok := err.(errors.IError); ok {
		if c := xerr.Cause(); c != nil {
			tab(log)
			log.Print(1, "Caused by:\n")
			log.XError(c, style, tab.Shift())
		}
	}
}

// Usage prints CLI usage information.
func (log Log) Usage() {
	fmt.Println("Please specify at least a spec file & an operation to test.")
	fmt.Println("Example:")
	fmt.Println("oasis from path/to/oas_spec.yaml test operation_id")
}

// LoadingSpec informs about the API specification being used.
func (log Log) LoadingSpec(path string) {
	log.Println(2, "Loading %s", log.Style.URL(path))
}

// PrintOperations prints the list of available operations.
func (log Log) PrintOperations(ops []*api.Operation) {
	for _, op := range ops {
		if op.ID != "" {
			log.Println(1, "\t%s [%s]", log.Style.Op(op.Name), log.Style.Op(op.ID))
			if op.Description != "" {
				log.Println(1, "\t%s", op.Description)
			}
		} else {
			log.Println(1, "\t%s", log.Style.Op(op.Name))
		}
		log.Println(1, "\t%s @ %s\n", op.Method, log.Style.URL(op.Path))
		log.Println(1, "")
	}
}

// TestingProject informs about the project being tested.
func (log Log) TestingProject(pi *api.ProjectInfo) {
	log.Println(2, "Testing the %s @ %s", log.Style.Op(pi.Title), log.Style.ID(pi.Version))
}

// UsingHost informs about the API host being used for testing.
func (log Log) UsingHost(host *api.Host) {
	log.Println(2, "Using the %s host @ %s", log.Style.Op(host.Name), log.Style.URL(host.URL))
}

// UsingDefaultHost informs that a default host has been chosen fr testing.
func (log Log) UsingDefaultHost() {
	log.Println(2, "No host name has been specified, using the first one in the list.")
}

// HostNotFound informs that no specified host has been found in the spec.
func (log Log) HostNotFound(h string) {
	if h == "" {
		log.Println(2, "No default host is found in the spec.")
	} else {
		log.Println(2, "The host \"%s\" is not found in the spec.", h)
	}
}

// UsingSecurity informs about security mechanisms being used during testing.
func (log Log) UsingSecurity(sec api.ISecurity) {
	log.Println(3, "\tUsing the %s security settings.", log.Style.ID(sec.GetName()))
}

// SecurityHasNoData informs that the selected security settings has no data to use in requests.
func (log Log) SecurityHasNoData(sec api.ISecurity) {
	log.Println(3, "\tThe security %s contains no data to use in request.", log.Style.ID(sec.GetName()))
}

// Requesting informs about an HTTP request being performed.
func (log Log) Requesting(method string, URL string) {
	log.Println(2, "\tRequesting %s @ %s", log.Style.Method(method), log.Style.URL(URL))
}

// ParameterHasNoExample informs that a parameter has no example value to use.
func (log Log) ParameterHasNoExample(paramName string, in string, container string) {
	log.Println(5, "\tThe %s parameter %s (from %s) has no example value to use.", in, log.Style.ID(paramName), container)
}

// UsingParameterExample informs that a parameter example being used.
func (log Log) UsingParameterExample(paramName string, in string, container string) {
	log.Println(5, "\tUsing the %s parameter %s (from %s) example.", in, log.Style.ID(paramName), container)
}

// HeaderHasNoValue informs that a required response header has no data.
func (log Log) HeaderHasNoValue(hdr *api.Header) {
	log.Println(1, "\tHeader \"%s\" is required but is not present.", hdr.Name)
}

// ResponseHasWrongStatus informs that the received response has wrong/unexpected status.
func (log Log) ResponseHasWrongStatus(resp *api.Response, actualStatus int) {
	m := strings.Join([]string{
		"\t",
		"Expected the %s ",
		log.Style.ID("status"),
		" in response, but got %s",
		".",
	}, "")
	log.Println(2, m, log.Style.ValueExpected(resp.StatusCode), log.Style.ValueActual(actualStatus))
}

// ResponseHasWrongContentType informs that the received response has wrong/unexpected Content-Type header value.
func (log Log) ResponseHasWrongContentType(resp *api.Response, actualCT string) {
	m := strings.Join([]string{
		"\t",
		"Expected the %s ",
		log.Style.ID("Content-Type"),
		" in response, but got %s",
		".",
	}, "")

	log.Println(2, m, log.Style.ValueExpected(resp.ContentType), log.Style.ValueActual(actualCT))
}

// TestingOperation informs about an operation being tested.
func (log Log) TestingOperation(op *api.Operation) {
	log.Print(1, "Testing the %s operation... ", log.Style.Op(op.Name))
	log.Print(2, "\n")
}

// OperationOK informs that the operation has finished successfully.
func (log Log) OperationOK(res *api.Operation) {
	log.Print(2, "\t")
	log.Println(1, "%s", log.Style.OK("SUCCESS"))
	log.Print(2, "\n")
}

// OperationFail informs that the operation has failed.
func (log Log) OperationFail(res *api.Operation) {
	log.Print(2, "\t")
	log.Println(1, "%s", log.Style.Failure("FAILURE"))
	log.Print(2, "\n")
}

// OperationNotFound informs that the requested operation is not found in the spec.
func (log Log) OperationNotFound(op string) {
	log.Println(1, "The operation \"%s\" isn't there.", op)
}

// SchemaTesting informs about a value being tested againt some JSON schema.
func (log Log) SchemaTesting(schema *api.Schema, data interface{}) {
	datas := log.Style.ValueActual(fmt.Sprintf("%#v", data))
	log.Print(4, "\t%s: testing %s", log.Style.ID(schema.Name), datas)
}

// SchemaOK informs that JSON schema testing finished successfully.
func (log Log) SchemaOK(schema *api.Schema) {
	log.Println(4, log.Style.Success(" - OK"))
}

// SchemaFail informs that JSON schema testing finished unsuccessfully.
func (log Log) SchemaFail(schema *api.Schema, errors []gojsonschema.ResultError) {
	log.Println(4, log.Style.Error(" - FAILURE"))

	for _, desc := range errors {
		log.Println(4, "\t\t%s", log.Style.Error(desc))
	}
}
