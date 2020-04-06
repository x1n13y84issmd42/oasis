package log

import (
	"fmt"
	"strings"

	"github.com/gookit/color"
	"github.com/x1n13y84issmd42/oasis/src/api"
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
	styleDefault       ColorFn
	styleURL           ColorFn
	styleMethod        ColorFn
	styleOp            ColorFn
	styleOK            ColorFn
	styleFailure       ColorFn
	styleSuccess       ColorFn
	styleError         ColorFn
	styleID            ColorFn
	styleValueExpected ColorFn
	styleValueActual   ColorFn
}

// NewFestive is a Nice logger constructor.
func NewFestive(level int64) *Festive {
	return &Festive{
		Log: Log{
			Level: level,
		},

		styleDefault:       wrapStyle(color.New(38, 5, 218).Sprint),
		styleURL:           wrapStyle(color.New(color.FgCyan, color.OpUnderscore).Sprint),
		styleMethod:        wrapStyle(color.New(color.FgCyan).Sprint),
		styleOp:            wrapStyle(color.New(color.FgYellow).Sprint),
		styleOK:            wrapStyle(color.New(color.FgLightWhite, color.BgGreen).Sprint),
		styleFailure:       wrapStyle(color.New(color.FgLightWhite, color.BgRed).Sprint),
		styleSuccess:       wrapStyle(color.New(color.FgGreen).Sprint),
		styleError:         wrapStyle(color.New(color.FgRed).Sprint),
		styleID:            wrapStyle(color.New(color.FgLightWhite, color.OpUnderscore).Sprint),
		styleValueExpected: wrapStyle(color.New(48, 5, 2, 38, 5, 0).Sprint),
		styleValueActual:   wrapStyle(color.New(48, 5, 240, 38, 5, 255).Sprint),
	}
}

// Usage prints CLI usage information.
func (log Festive) Usage() {
	fmt.Println("Please specify at least a spec file & an operation to test.")
	fmt.Println("Example:")
	fmt.Println("oasis from path/to/oas_spec.yaml test operation_id")
}

// Error --
func (log Festive) Error(err error) {
	log.Println(1, "\tSomething happened: %s", log.styleError(err.Error()))
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

// Overriding --
func (log Festive) Overriding(what string) {
	log.Println(1, "\tOverriding %s.", what)
}

// Requesting --
func (log Festive) Requesting(method string, URL string) {
	log.Println(2, "\tRequesting %s @ %s", log.styleMethod(method), log.styleURL(URL))
}

// ResponseNotFound --
func (log Festive) ResponseNotFound(CT string, status int) {
	log.Println(1, "\tNo response for Status of %d & Content-Type of \"%s\"", status, CT)
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

// UsingRequest --
/* func (log Nice) UsingRequest(req *api.Request) {
	log.Println(1, "\tUsing the \"%s\" request.", req.ContentType)
} */

// UsingResponse --
func (log Festive) UsingResponse(resp *api.Response) {
	// if resp.Schema != nil {
	// 	log.Println(1, "\tTesting against the \"%s\" response.", resp.Schema.Name)
	// } else {
	CT := resp.ContentType
	if len(CT) == 0 {
		CT = "*/*"
	}
	log.Println(1, "\tTesting against the %s @ %d response.", CT, resp.StatusCode)
	// }
}

// HeaderHasNoValue --
func (log Festive) HeaderHasNoValue(hdr *api.Header) {
	log.Println(1, "\tHeader \"%s\" is required but is not present.", hdr.Name)
}

// HeaderHasWrongType --
/* func (log Nice) HeaderHasWrongType(hdr *api.Header) {
	log.Println(1, "\tHeader \"%s\" has a wrong type.", hdr.Name)
} */

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

// UnknownSchemaDataType --
/* func (log Nice) UnknownSchemaDataType(schema *api.Schema) {
	log.Println(1, "\tSchema \"%s\" has unknown data type \"%s\".", schema.Name, schema.DataType)
} */

// SchemaExpectedBoolean --
/* func (log Nice) SchemaExpectedBoolean(schema *api.Schema, v interface{}) {
	log.Println(1, "\tSchema \"%s\" expected %#v to be a boolean type.", schema.Name, v)
} */

// SchemaExpectedNumber --
/* func (log Nice) SchemaExpectedNumber(schema *api.Schema, v interface{}) {
	log.Println(1, "\tSchema \"%s\" expected %#v to be a floating point number.", schema.Name, v)
} */

// SchemaExpectedInteger --
/* func (log Nice) SchemaExpectedInteger(schema *api.Schema, v interface{}) {
	log.Println(1, "\tSchema \"%s\" expected %#v to be an integer number.", schema.Name, v)
} */

// SchemaExpectedString --
/* func (log Nice) SchemaExpectedString(schema *api.Schema, v interface{}) {
	log.Println(1, "\tSchema \"%s\" expected %#v to be a string type.", schema.Name, v)
} */

// SchemaExpectedArray --
/* func (log Nice) SchemaExpectedArray(schema *api.Schema, v interface{}) {
	log.Println(1, "\tSchema \"%s\" expected %#v to be an array type.", schema.Name, v)
} */

// SchemaExpectedObject --
/* func (log Nice) SchemaExpectedObject(schema *api.Schema, v interface{}) {
	log.Println(1, "\tSchema \"%s\" expected %#v to be an object type.", schema.Name, v)
} */

// UsingSecurity --
/* func (log Nice) UsingSecurity(sec *api.Security) {
	log.Println(1, "\tUsing the \"%s\" security settings.", sec.Name)
} */

// ParameterHasNoExample --
func (log Festive) ParameterHasNoExample(paramName string, in string, container string) {
	log.Println(5, "\tThe %s parameter \"%s\" (from %s) has no example value to use.", in, paramName, container)
}

// UsingParameterExample --
func (log Festive) UsingParameterExample(paramName string, in string, container string) {
	log.Println(5, "\tUsing the %s parameter \"%s\" (from %s) example.", in, paramName, container)
}

// PropertyHasNoValue --
/* func (log Nice) PropertyHasNoValue(prop *api.Property, ctx *utility.Context) {
	log.Println(1, "\t%s: property is required but is not present.", ctx.String())
} */

// PropertyHasWrongType --
/* func (log Nice) PropertyHasWrongType(prop *api.Property, ctx *utility.Context) {
	log.Println(1, "\t%s: property has wrong type. Expected %s, got %s.", ctx.String(), prop.Schema.DataType, ctx.CurrentValueType())
} */

// TestingProject --
func (log Festive) TestingProject(pi *api.ProjectInfo) {
	log.Println(1, "Testing the %s @ %s", log.styleOp(pi.Title), log.styleID(pi.Version))
}

// UsingHost --
func (log Festive) UsingHost(host *api.Host) {
	log.Println(2, "Using the %s host @ %s", log.styleOp(host.Name), log.styleURL(host.URL))
}
