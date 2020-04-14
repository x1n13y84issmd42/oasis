package log

import (
	"fmt"
	"strings"

	"github.com/gookit/color"
	"github.com/x1n13y84issmd42/oasis/src/api"
	"github.com/x1n13y84issmd42/oasis/src/errors"
	"github.com/xeipuuv/gojsonschema"
)

// ColorFn is a function to colorize strings before printing them.
type ColorFn = func(...interface{}) string

func wrapStyle(styleFn ColorFn) ColorFn {
	return func(args ...interface{}) string {
		// Adding a "redraw the prompt" control sequence.
		// This is needed for new BG color to take place in the mid of the line.
		// Otherwise it redraws only on \n.
		return styleFn(args...) + "\x1b[K"
	}
}

// Festive - a colorized test execution logger.
type Festive struct {
	Log
}

// NewFestive is a Nice logger constructor.
func NewFestive(level int64) *Festive {
	return &Festive{
		Log: Log{
			Level: level,
		},
	}
}

func (log Festive) styleDefault(args ...interface{}) string {
	return color.New(38, 5, 218).Sprint(args...) + "\x1b[K"
}

func (log Festive) styleURL(args ...interface{}) string {
	return color.New(color.FgCyan, color.OpUnderscore).Sprint(args...) + "\x1b[K"
}

func (log Festive) styleMethod(args ...interface{}) string {
	return color.New(color.FgCyan).Sprint(args...) + "\x1b[K"
}

func (log Festive) styleOp(args ...interface{}) string {
	return color.New(color.FgYellow).Sprint(args...) + "\x1b[K"
}

func (log Festive) styleOK(args ...interface{}) string {
	return color.New(color.FgLightWhite, color.BgGreen).Sprint(args...) + "\x1b[K"
}

func (log Festive) styleFailure(args ...interface{}) string {
	return color.New(color.FgLightWhite, color.BgRed).Sprint(args...) + "\x1b[K"
}

func (log Festive) styleSuccess(args ...interface{}) string {
	return color.New(color.FgGreen).Sprint(args...) + "\x1b[K"
}

func (log Festive) styleError(args ...interface{}) string {
	return color.New(color.FgRed).Sprint(args...) + "\x1b[K"
}

func (log Festive) styleID(args ...interface{}) string {
	return color.New(color.FgLightWhite, color.OpUnderscore).Sprint(args...) + "\x1b[K"
}

func (log Festive) styleValueExpected(args ...interface{}) string {
	return color.New(48, 5, 2, 38, 5, 0).Sprint(args...) + "\x1b[K"
}

func (log Festive) styleValueActual(args ...interface{}) string {
	return color.New(48, 5, 240, 38, 5, 255).Sprint(args...) + "\x1b[K"
}

// Usage prints CLI usage information.
func (log Festive) Usage() {
	fmt.Println("Please specify at least a spec file & an operation to test.")
	fmt.Println("Example:")
	fmt.Println("oasis from path/to/oas_spec.yaml test operation_id")
}

// Error --
func (log Festive) Error(err error) {
	if xerr, ok := err.(errors.IError); ok {
		// log.Println(1, "MEGA ERROR: %s", log.styleError(err.Error()))
		log.XError(xerr, log, Tab(0))
		return
	}
	// log.Println(1, "\tSomething happened: %s", log.styleError(err.Error()))
	log.Println(1, log.styleError(err.Error()))
}

// LoadingSpec --
func (log Festive) LoadingSpec(path string) {
	log.Println(2, "Loading %s", log.styleURL(path))
}

// PrintOperations prints the list of available operations.
func (log Festive) PrintOperations(ops []*api.Operation) {
	for _, op := range ops {
		if op.ID != "" {
			log.Println(1, "\t%s [%s]", log.styleOp(op.Name), log.styleOp(op.ID))
			if op.Description != "" {
				log.Println(1, "\t%s", op.Description)
			}
		} else {
			log.Println(1, "\t%s", log.styleOp(op.Name))
		}
		log.Println(1, "\t%s @ %s\n", op.Method, log.styleURL(op.Path))
		log.Println(1, "")
	}
}

// TestingProject --
func (log Festive) TestingProject(pi *api.ProjectInfo) {
	log.Println(2, "Testing the %s @ %s", log.styleOp(pi.Title), log.styleID(pi.Version))
}

// UsingHost --
func (log Festive) UsingHost(host *api.Host) {
	log.Println(2, "Using the %s host @ %s", log.styleOp(host.Name), log.styleURL(host.URL))
}

// UsingDefaultHost --
func (log Festive) UsingDefaultHost() {
	log.Println(2, "No host name has been specified, using the first one in the list.")
}

// HostNotFound ...
func (log Festive) HostNotFound(h string) {
	if h == "" {
		log.Println(2, "No default host is found in the spec.")
	} else {
		log.Println(2, "The host \"%s\" is not found in the spec.", h)
	}
}

// UsingSecurity --
func (log Festive) UsingSecurity(sec api.ISecurity) {
	log.Println(3, "\tUsing the %s security settings.", log.styleID(sec.GetName()))
}

// SecurityHasNoData --
func (log Festive) SecurityHasNoData(sec api.ISecurity) {
	log.Println(3, "\tThe security %s contains no data to use in request.", log.styleID(sec.GetName()))
}

// Requesting --
func (log Festive) Requesting(method string, URL string) {
	log.Println(2, "\tRequesting %s @ %s", log.styleMethod(method), log.styleURL(URL))
}

// ParameterHasNoExample --
func (log Festive) ParameterHasNoExample(paramName string, in string, container string) {
	log.Println(5, "\tThe %s parameter %s (from %s) has no example value to use.", in, log.styleID(paramName), container)
}

// UsingParameterExample --
func (log Festive) UsingParameterExample(paramName string, in string, container string) {
	log.Println(5, "\tUsing the %s parameter %s (from %s) example.", in, log.styleID(paramName), container)
}

// HeaderHasNoValue --
func (log Festive) HeaderHasNoValue(hdr *api.Header) {
	log.Println(1, "\tHeader \"%s\" is required but is not present.", hdr.Name)
}

// ResponseHasWrongStatus --
func (log Festive) ResponseHasWrongStatus(resp *api.Response, actualStatus int) {
	m := strings.Join([]string{
		"\t",
		"Expected the %s ",
		log.styleID("status"),
		" in response, but got %s",
		".",
	}, "")
	log.Println(2, m, log.styleValueExpected(resp.StatusCode), log.styleValueActual(actualStatus))
}

// ResponseHasWrongContentType --
func (log Festive) ResponseHasWrongContentType(resp *api.Response, actualCT string) {
	m := strings.Join([]string{
		"\t",
		"Expected the %s ",
		log.styleID("Content-Type"),
		" in response, but got %s",
		".",
	}, "")

	log.Println(2, m, log.styleValueExpected(resp.ContentType), log.styleValueActual(actualCT))
}

// TestingOperation --
func (log Festive) TestingOperation(op *api.Operation) {
	log.Print(1, "Testing the %s operation... ", log.styleOp(op.Name))
	log.Print(2, "\n")
}

// OperationOK --
func (log Festive) OperationOK(res *api.Operation) {
	log.Print(2, "\t")
	log.Println(1, "%s", log.styleOK("SUCCESS"))
	log.Print(2, "\n")
}

// OperationFail --
func (log Festive) OperationFail(res *api.Operation) {
	log.Print(2, "\t")
	log.Println(1, "%s", log.styleFailure("FAILURE"))
	log.Print(2, "\n")
}

// OperationNotFound --
func (log Festive) OperationNotFound(op string) {
	log.Println(1, "The operation \"%s\" isn't there.", op)
}

// SchemaTesting --
func (log Festive) SchemaTesting(schema *api.Schema, data interface{}) {
	datas := log.styleValueActual(fmt.Sprintf("%#v", data))
	log.Print(4, "\t%s: testing %s", log.styleID(schema.Name), datas)
}

// SchemaOK --
func (log Festive) SchemaOK(schema *api.Schema) {
	log.Println(4, log.styleSuccess(" - OK"))
}

// SchemaFail --
func (log Festive) SchemaFail(schema *api.Schema, errors []gojsonschema.ResultError) {
	log.Println(4, log.styleError(" - FAILURE"))
	// log.Println(4, "\tSchema \"%s\" has errors.", schema.Name)

	for _, desc := range errors {
		log.Println(4, "\t\t%s", log.styleError(desc))
	}
}

// ErrOperationMalformed ...
func (log Festive) ErrOperationMalformed(err *api.ErrOperationMalformed) {
	log.Println(1, "Operation %s has malformed or incomplete data.", log.styleID(err.OpID))
}

// ErrOperationNotFound ...
func (log Festive) ErrOperationNotFound(err *api.ErrOperationNotFound) {
	log.Println(1, "Operation %s not found.", log.styleID(err.OpID))
}
