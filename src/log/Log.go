package log

import (
	"fmt"
	"os"
	"strings"

	"github.com/x1n13y84issmd42/oasis/src/api"
	"github.com/x1n13y84issmd42/oasis/src/contract"
	"github.com/x1n13y84issmd42/oasis/src/errors"
	"github.com/xeipuuv/gojsonschema"
)

// Log is a base type for loggers.
type Log struct {
	Level int64
	Style contract.LogStyle
}

// New creates a new logger based on the provided log style & level.
func New(style string, level int64) contract.Logger {
	switch style {
	case "plain":
		return NewPlain(level)

	case "festive":
		return NewFestive(level)
	}

	fmt.Printf("The \"%s\" log style is unknown.\nAvailable loggers are:\n", style)
	fmt.Println("\tplain - a plain text logger")
	fmt.Println("\tfestive - a nicer colorized logger")

	os.Exit(1)

	return nil
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

// Parameters prints a ParameterSource contents.
func (log Log) Parameters(name string, params contract.ParameterSource) {
	log.Println(0, "Contents of %s", name)
	for p := range params.Iterate() {
		fmt.Printf("%s = %s (from %s)\n", log.Style.ID(p.N), log.Style.ValueActual(p.V()), p.Source)
	}
	log.Println(0, "")
}

// Error outputs errors. It accepts both built-in errors as well as errors.IError instances.
// The latter carry their cause error, if any, and those will be recursively printed.
func (log Log) Error(err error) {
	log.XError(err, log.Style, contract.Tab(0))
	log.Println(1, "")
}

// XError is an internal error handling function. It handles both
// built-in errors and errors.IError instances.
func (log Log) XError(err error, style contract.LogStyle, tab contract.TabFn) {
	tab(log)
	log.Println(1, "%s", style.Error(err.Error()))

	if xerr, ok := err.(errors.IError); ok {

		log.Println(1, "")
		log.Println(1, "%s", style.Error("Source: "+xerr.Caller()))

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
	log.Println(2, "Loading the %s spec.", log.Style.URL(path))
}

// LoadingScript informs about the API specification being used.
func (log Log) LoadingScript(path string) {
	log.Println(2, "Loading the %s script.", log.Style.URL(path))
}

// PrintOperations prints the list of available operations.
func (log Log) PrintOperations(ops contract.OperationIterator) {
	for op := range ops {
		if op.ID() != "" {
			log.Println(1, "\t%s [%s]", log.Style.Op(op.Name()), log.Style.Op(op.ID()))
			if op.Description() != "" {
				log.Println(1, "\t%s", strings.ReplaceAll(op.Description(), "\n", "\n\t"))
			}
		} else {
			log.Println(1, "\t%s", log.Style.Op(op.Name()))
		}
		log.Println(1, "\t%s @ %s\n", op.Method(), log.Style.URL(op.Path()))
		// log.Println(1, "")
	}
}

// TestingProject informs about the project being tested.
func (log Log) TestingProject(pi contract.ProjectInfo) {
	log.Println(2, "Testing the %s @ %s", log.Style.Op(pi.Title()), log.Style.ID(pi.Version()))
}

// UsingSecurity informs about security mechanisms being used during testing.
func (log Log) UsingSecurity(sec contract.Security) {
	log.Println(3, "\tUsing the %s security settings.", log.Style.ID(sec.GetName()))
}

// SecurityHasNoData informs that the selected security settings has no data to use in requests.
func (log Log) SecurityHasNoData(sec contract.Security) {
	log.Println(3, "\tThe security %s contains no data to use in request.", log.Style.ID(sec.GetName()))
}

// Requesting informs about an HTTP request being performed.
func (log Log) Requesting(method string, URL string) {
	log.Println(2, "\tRequesting %s @ %s", log.Style.Method(method), log.Style.URL(URL))
}

// UsingParameterExample informs that a parameter example being used.
func (log Log) UsingParameterExample(paramName string, in string, container string, value string) {
	log.Println(5, "\tUsing the %s parameter %s %s (from %s).", in, log.Style.ID(paramName), log.Style.ValueActual(value), container)
}

// Expecting informs that a parameter example being used.
func (log Log) Expecting(what string, v string) {
	log.Println(5, "\tExpecting %s %s.", log.Style.ID(what), log.Style.ValueActual(v))
}

// HeaderHasNoValue informs that a required response header has no data.
func (log Log) HeaderHasNoValue(hdr string) {
	log.Println(1, "\tHeader \"%s\" is required but is not present.", hdr)
}

// ResponseHasWrongStatus informs that the received response has wrong/unexpected status.
func (log Log) ResponseHasWrongStatus(expectedStatus int, actualStatus int) {
	m := strings.Join([]string{
		"\t",
		"Expected the %s ",
		log.Style.ID("status"),
		" in response, but got %s",
		".",
	}, "")
	log.Println(2, m, log.Style.ValueExpected(expectedStatus), log.Style.ValueActual(actualStatus))
}

// ResponseHasWrongContentType informs that the received response has wrong/unexpected Content-Type header value.
func (log Log) ResponseHasWrongContentType(expectedCT string, actualCT string) {
	m := strings.Join([]string{
		"\t",
		"Expected the %s ",
		log.Style.ID("Content-Type"),
		" in response, but got %s",
		".",
	}, "")

	log.Println(2, m, log.Style.ValueExpected(expectedCT), log.Style.ValueActual(actualCT))
}

// TestingOperation informs about an operation being tested.
func (log Log) TestingOperation(op contract.Operation) {
	log.Print(1, "Testing the %s operation... ", log.Style.Op(op.Name()))
	log.Print(2, "\n")
}

// OperationOK informs that the operation has finished successfully.
func (log Log) OperationOK() {
	log.Print(2, "\t")
	log.Println(1, "%s", log.Style.OK("SUCCESS"))
	log.Print(2, "\n")
}

// OperationFail informs that the operation has failed.
func (log Log) OperationFail() {
	log.Print(2, "\t")
	log.Println(1, "%s", log.Style.Failure("FAILURE"))
	log.Print(2, "\n")
}

// SchemaTesting informs about a value being tested againt some JSON schema.
func (log Log) SchemaTesting(schema *api.Schema, data interface{}) {
	datas := log.Style.ValueActual(fmt.Sprintf("%#v", data))
	log.Print(4, "\t%s: testing %s", log.Style.ID(schema.Name), datas)
}

// SchemaOK informs that JSON schema testing finished successfully.
func (log Log) SchemaOK() {
	log.Println(3, "\t%s", log.Style.Success("Schema OK."))
}

// SchemaFail informs that JSON schema testing finished unsuccessfully.
func (log Log) SchemaFail(errors []gojsonschema.ResultError) {
	log.Println(3, "\t%s", log.Style.Error("Schema failure."))

	for _, desc := range errors {
		log.Println(3, "\t%s", log.Style.Error(desc))
	}
}

// ScriptExecutionStart logs the starting node of the script execution graph.
func (log Log) ScriptExecutionStart(node string) {
	log.Println(3, "Execution starts from the node %s.\n", log.Style.Op(node))
}
