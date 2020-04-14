package log

import (
	"fmt"

	"github.com/x1n13y84issmd42/oasis/src/api"
	"github.com/xeipuuv/gojsonschema"
)

// Plain - a plain & simple test execution logger
type Plain struct {
	Log
}

// NewPlain is a Nice logger constructor.
func NewPlain(level int64) *Plain {
	return &Plain{
		Log: Log{
			Level: level,
		},
	}
}

func (log Plain) styleDefault(args ...interface{}) string {
	return fmt.Sprint(args...)
}

func (log Plain) styleURL(args ...interface{}) string {
	return fmt.Sprint(args...)
}

func (log Plain) styleMethod(args ...interface{}) string {
	return fmt.Sprint(args...)
}

func (log Plain) styleOp(args ...interface{}) string {
	return fmt.Sprint(args...)
}

func (log Plain) styleOK(args ...interface{}) string {
	return fmt.Sprint(args...)
}

func (log Plain) styleFailure(args ...interface{}) string {
	return fmt.Sprint(args...)
}

func (log Plain) styleSuccess(args ...interface{}) string {
	return fmt.Sprint(args...)
}

func (log Plain) styleError(args ...interface{}) string {
	return fmt.Sprint(args...)
}

func (log Plain) styleID(args ...interface{}) string {
	return fmt.Sprint(args...)
}

func (log Plain) styleValueExpected(args ...interface{}) string {
	return fmt.Sprint(args...)
}

func (log Plain) styleValueActual(args ...interface{}) string {
	return fmt.Sprint(args...)
}

// Usage prints CLI usage information.
func (log Plain) Usage() {
	fmt.Println("Please specify at least a spec file & an operation to test.")
	fmt.Println("Example:")
	fmt.Println("oasis from path/to/oas_spec.yaml test operation_id")
}

// Error --
func (log Plain) Error(err error) {
	log.XError(err, log, Tab(0))
}

// LoadingSpec --
func (log Plain) LoadingSpec(path string) {
	log.Println(2, "Loading %s", path)
}

// PrintOperations prints the list of available operations.
func (log Plain) PrintOperations(ops []*api.Operation) {
	for _, op := range ops {
		if op.ID != "" {
			log.Println(1, "\t%s [%s]", op.Name, op.ID)
			if op.Description != "" {
				log.Println(1, "\t%s", op.Description)
			}
		} else {
			log.Println(1, "\t%s", op.Name)
		}
		log.Println(1, "\t%s @ %s\n", op.Method, op.Path)
		log.Println(1, "")
	}
}

// TestingProject --
func (log Plain) TestingProject(pi *api.ProjectInfo) {
	log.Println(1, "Testing the %s @ %s", pi.Title, pi.Version)
}

// UsingHost --
func (log Plain) UsingHost(host *api.Host) {
	log.Println(2, "Using the \"%s\" host @ %s", host.Name, host.URL)
}

// UsingDefaultHost --
func (log Plain) UsingDefaultHost() {
	log.Println(2, "No host name has been specified, using the first one in the list.")
}

// HostNotFound ...
func (log Plain) HostNotFound(h string) {
	if h == "" {
		log.Println(2, "No default host is found in the spec.")
	} else {
		log.Println(2, "The host \"%s\" is not found in the spec.", h)
	}
}

// UsingSecurity --
func (log Plain) UsingSecurity(sec api.ISecurity) {
	log.Println(3, "\tUsing the \"%s\" security settings.", sec.GetName())
}

// SecurityHasNoData --
func (log Plain) SecurityHasNoData(sec api.ISecurity) {
	log.Println(3, "\tThe security %s contains no data to use in request.", sec.GetName())
}

// Requesting --
func (log Plain) Requesting(method string, URL string) {
	log.Println(2, "\tRequesting %s @ %s", method, URL)
}

// ParameterHasNoExample --
func (log Plain) ParameterHasNoExample(paramName string, in string, container string) {
	log.Println(5, "\tThe %s parameter \"%s\" (from %s) has no example value to use.", in, paramName, container)
}

// UsingParameterExample --
func (log Plain) UsingParameterExample(paramName string, in string, container string) {
	log.Println(5, "\tUsing the %s parameter \"%s\" (from %s) example.", in, paramName, container)
}

// HeaderHasNoValue --
func (log Plain) HeaderHasNoValue(hdr *api.Header) {
	log.Println(1, "\tHeader \"%s\" is required but is not present.", hdr.Name)
}

// ResponseHasWrongStatus --
func (log Plain) ResponseHasWrongStatus(resp *api.Response, actualStatus int) {
	log.Println(1, "\tExpected the %d status in response, but got %d.", resp.StatusCode, actualStatus)
}

// ResponseHasWrongContentType --
func (log Plain) ResponseHasWrongContentType(resp *api.Response, actualCT string) {
	log.Println(1, "\tExpected the \"%s\" Content-Type in response, but got \"%s\".", resp.ContentType, actualCT)
}

// TestingOperation --
func (log Plain) TestingOperation(res *api.Operation) {
	log.Print(1, "Testing the \"%s\" operation... ", res.Name)
	log.Print(2, "\n")
}

// OperationOK --
func (log Plain) OperationOK(res *api.Operation) {
	log.Print(2, "\t")
	log.Println(1, "%s", "SUCCESS")
	log.Print(2, "\n")
}

// OperationFail --
func (log Plain) OperationFail(res *api.Operation) {
	log.Print(2, "\t")
	log.Println(1, "%s", "FAILURE")
	log.Print(2, "\n")
}

// OperationNotFound --
func (log Plain) OperationNotFound(op string) {
	log.Println(1, "The operation \"%s\" isn't there.", op)
}

// SchemaTesting --
func (log Plain) SchemaTesting(schema *api.Schema, data interface{}) {
	log.Print(4, "\t%s: testing %#v", schema.Name, data)
}

// SchemaOK --
func (log Plain) SchemaOK(schema *api.Schema) {
	log.Println(4, " - OK")
}

// SchemaFail --
func (log Plain) SchemaFail(schema *api.Schema, errors []gojsonschema.ResultError) {
	log.Println(4, " - FAILURE")
	// log.Println(4, "\tSchema \"%s\" has errors.", schema.Name)

	for _, desc := range errors {
		log.Println(4, "\t\t%s", desc)
	}
}

// ErrOperationMalformed ...
func (log Plain) ErrOperationMalformed(err *api.ErrOperationMalformed) {
	log.Println(1, "Oops")
}

// ErrOperationNotFound ...
func (log Plain) ErrOperationNotFound(err *api.ErrOperationNotFound) {
	log.Println(1, "Oops")
}
