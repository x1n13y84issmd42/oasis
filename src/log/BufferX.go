package log

import (
	"fmt"
	"strings"
	"sync"

	"github.com/x1n13y84issmd42/oasis/src/api"
	"github.com/x1n13y84issmd42/oasis/src/contract"
	"github.com/x1n13y84issmd42/oasis/src/errors"
	"github.com/xeipuuv/gojsonschema"
)

// BufferX accumulates log output to flush it all at once.
// Useful for concurrent script execution.
type BufferX struct {
	Log
	data  string
	Mutex sync.Mutex
}

// NewBufferX creates a new BufferX instance.
func NewBufferX(logger *Log) contract.Logger {
	return &BufferX{
		Log: *logger,
	}
}

// Print prints.
func (log *BufferX) Print(l int64, msg string, args ...interface{}) {
	if l <= log.Level {
		data := fmt.Sprintf(msg, args...)
		// fmt.Printf("Captured %s\n", data)
		log.data += data
	}
}

// Println prints and adds a newline.
func (log *BufferX) Println(l int64, msg string, args ...interface{}) {
	log.Print(l, msg+"\n", args...)
}

// Flush flushes the accumulated output to stdout.
func (log *BufferX) Flush() {
	// fmt.Print("Flushing\n")
	log.Mutex.Lock()
	fmt.Print(log.data)
	log.data = ""
	log.Mutex.Unlock()
}

// NOMESSAGE is a default and temporary print function to use when you don't have
// a dedidated message function in the logger. It is meant to be replaced eventually
// with a proper logging method or removed. Henec the indicating name.
func (log *BufferX) NOMESSAGE(msg string, args ...interface{}) {
	log.Println(1, "\t"+msg, args...)
}

// Parameters prints a ParameterSource contents.
func (log *BufferX) Parameters(name string, params contract.ParameterSource) {
	log.Println(0, "Contents of %s", name)
	for p := range params.Iterate() {
		fmt.Printf("%s = %s (from %s)\n", log.Style.ID(p.N), log.Style.Value(p.V()), p.Source)
	}
	log.Println(0, "")
}

// Error outputs errors. It accepts both built-in errors as well as errors.IError instances.
// The latter carry their cause error, if any, and those will be recursively printed.
func (log *BufferX) Error(err error) {
	log.XError(err, log.Style, contract.Tab(0))
	log.Println(1, "")
}

// XError is an internal error handling function. It handles both
// built-in errors and errors.IError instances.
func (log *BufferX) XError(err error, style contract.LogStyle, tab contract.TabFn) {
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
func (log *BufferX) Usage() {
	fmt.Println("Please specify at least a spec file & an operation to test.")
	fmt.Println("Example:")
	fmt.Println("oasis from path/to/oas_spec.yaml test operation_id")
}

// LoadingSpec informs about the API specification being used.
func (log *BufferX) LoadingSpec(path string) {
	log.Println(2, "Loading the %s spec.", log.Style.URL(path))
}

// LoadingScript informs about the API specification being used.
func (log *BufferX) LoadingScript(path string) {
	log.Println(2, "Loading the %s script.", log.Style.URL(path))
}

// PrintOperations prints the list of available operations.
func (log *BufferX) PrintOperations(ops contract.OperationIterator) {
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
func (log *BufferX) TestingProject(pi contract.ProjectInfo) {
	log.Println(2, "Testing the %s @ %s", log.Style.Op(pi.Title()), log.Style.ID(pi.Version()))
}

// UsingSecurity informs about security mechanisms being used during testing.
func (log *BufferX) UsingSecurity(sec contract.Security) {
	log.Println(3, "\tUsing the %s security settings.", log.Style.ID(sec.GetName()))
}

// SecurityHasNoData informs that the selected security settings has no data to use in requests.
func (log *BufferX) SecurityHasNoData(sec contract.Security) {
	log.Println(3, "\tThe security %s contains no data to use in request.", log.Style.ID(sec.GetName()))
}

// Requesting informs about an HTTP request being performed.
func (log *BufferX) Requesting(method string, URL string) {
	log.Println(2, "\tRequesting %s @ %s", log.Style.Method(method), log.Style.URL(URL))
}

// UsingParameterExample informs that a parameter example being used.
func (log *BufferX) UsingParameterExample(paramName string, in string, container string, value string) {
	log.Println(5, "\tUsing the %s parameter %s %s (from %s).", in, log.Style.ID(paramName), log.Style.Value(value), container)
}

// Expecting informs that a parameter example being used.
func (log *BufferX) Expecting(what string, v string) {
	log.Println(5, "\tExpecting %s %s.", log.Style.ID(what), log.Style.Value(v))
}

// ExpectingProperty informs that a parameter example being used.
func (log *BufferX) ExpectingProperty(what string, v string) {
	log.Println(5, "\tExpecting %s body property %s.", log.Style.ID(what), log.Style.Value(v))
}

// HeaderHasNoValue informs that a required response header has no data.
func (log *BufferX) HeaderHasNoValue(hdr string) {
	log.Println(1, "\tHeader \"%s\" is required but is not present.", hdr)
}

// ResponseHasWrongStatus informs that the received response has wrong/unexpected status.
func (log *BufferX) ResponseHasWrongStatus(expectedStatus int, actualStatus int) {
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
func (log *BufferX) ResponseHasWrongContentType(expectedCT string, actualCT string) {
	m := strings.Join([]string{
		"\t",
		"Expected the %s ",
		log.Style.ID("Content-Type"),
		" in response, but got %s",
		".",
	}, "")

	log.Println(2, m, log.Style.ValueExpected(expectedCT), log.Style.ValueActual(actualCT))
}

// ResponseHasWrongPropertyValue informs that the received response has wrong/unexpected body property value.
func (log *BufferX) ResponseHasWrongPropertyValue(propName string, expected string, actual string) {
	m := strings.Join([]string{
		"\t",
		"Expected the %s property to equal %s ",
		"but got %s",
		".",
	}, "")

	log.Println(2, m, log.Style.ID(propName), log.Style.ValueExpected(expected), log.Style.ValueActual(actual))
}

// TestingOperation informs about an operation being tested.
func (log *BufferX) TestingOperation(op contract.Operation) {
	log.Print(1, "Testing the %s operation... ", log.Style.Op(op.Name()))
	log.Print(2, "\n")
}

// OperationOK informs that the operation has finished successfully.
func (log *BufferX) OperationOK() {
	log.Print(2, "\t")
	log.Println(1, "%s", log.Style.OK("SUCCESS"))
	log.Print(2, "\n")
}

// OperationFail informs that the operation has failed.
func (log *BufferX) OperationFail() {
	log.Print(2, "\t")
	log.Println(1, "%s", log.Style.Failure("FAILURE"))
	log.Print(2, "\n")
}

// SchemaTesting informs about a value being tested againt some JSON schema.
func (log *BufferX) SchemaTesting(schema *api.Schema, data interface{}) {
	datas := log.Style.Value(fmt.Sprintf("%#v", data))
	log.Print(4, "\t%s: testing %s", log.Style.ID(schema.Name), datas)
}

// SchemaOK informs that JSON schema testing finished successfully.
func (log *BufferX) SchemaOK(schemaName string) {
	log.Println(4, "\t%s", log.Style.Success(schemaName+" schema OK."))
}

// SchemaFail informs that JSON schema testing finished unsuccessfully.
func (log *BufferX) SchemaFail(schemaName string, errors []gojsonschema.ResultError) {
	log.Println(4, "\t%s", log.Style.Error(schemaName+" schema failure."))

	for _, desc := range errors {
		log.Println(4, "\t%s", log.Style.Error(desc))
	}
}

// ScriptExecutionStart logs the starting node of the script execution graph.
func (log *BufferX) ScriptExecutionStart(node string) {
	log.Println(5, "Execution starts from the node %s.\n", log.Style.Op(node))
}
