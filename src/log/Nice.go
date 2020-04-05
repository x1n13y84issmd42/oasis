package log

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/x1n13y84issmd42/oasis/src/api"
	"github.com/xeipuuv/gojsonschema"
)

// ColorFn is a function to colorize strings before printing them.
type ColorFn = func(...interface{}) string

// Nice - a colorized test execution logger.
type Nice struct {
	Log
	colorURL     ColorFn
	colorOp      ColorFn
	colorOK      ColorFn
	colorFailure ColorFn
	colorSuccess ColorFn
	colorError   ColorFn
	colorID      ColorFn
	colorValue   ColorFn
}

// NewNice is a Nice logger constructor.
func NewNice(level int64) *Nice {
	return &Nice{
		Log: Log{
			Level: level,
		},

		colorURL:     color.New(color.FgCyan).Add(color.Underline).SprintFunc(),
		colorOp:      color.New(color.FgYellow).SprintFunc(),
		colorOK:      color.New(color.FgWhite).Add(color.BgGreen).SprintFunc(),
		colorFailure: color.New(color.FgWhite).Add(color.BgRed).SprintFunc(),
		colorSuccess: color.New(color.FgGreen).SprintFunc(),
		colorError:   color.New(color.FgRed).SprintFunc(),
		colorID:      color.New(color.FgHiWhite).SprintFunc(),
		colorValue:   color.New(color.FgHiWhite).SprintFunc(),
	}
}

// Usage prints CLI usage information.
func (log Nice) Usage() {
	fmt.Println("Please specify at least a spec file & an operation to test.")
	fmt.Println("Example:")
	fmt.Println("oasis from path/to/oas_spec.yaml test operation_id")
}

// Error --
func (log Nice) Error(err error) {
	log.Println(1, "\tSomething happened: %s", log.colorError(err.Error()))
}

// LoadingSpec --
func (log Nice) LoadingSpec(path string) {
	log.Println(2, "Loading %s", path)
}

// PrintOperations prints the list of available operations.
func (log Nice) PrintOperations(ops []*api.Operation) {
	for _, op := range ops {
		if op.ID != "" {
			log.Println(1, "\t%s [%s]", log.colorOp(op.Name), log.colorOp(op.ID))
			if op.Description != "" {
				log.Println(1, "\t%s", op.Description)
			}
		} else {
			log.Println(1, "\t%s", log.colorOp(op.Name))
		}
		log.Println(1, "\t%s @ %s\n", op.Method, log.colorURL(op.Path))
		log.Println(1, "")
	}
}

// UsingDefaultHost --
func (log Nice) UsingDefaultHost() {
	log.Println(2, "No host name has been specified, using the first one in the list.")
}

func (log Nice) HostNotFound(h string) {
	if h == "" {
		log.Println(2, "No default host is found in the spec.")
	} else {
		log.Println(2, "The host \"%s\" is not found in the spec.", h)
	}
}

// Overriding --
func (log Nice) Overriding(what string) {
	log.Println(1, "\tOverriding %s.", what)
}

// Requesting --
func (log Nice) Requesting(URL string) {
	log.Println(2, "\tRequesting %s", log.colorURL(URL))
}

// ResponseNotFound --
func (log Nice) ResponseNotFound(CT string, status int) {
	log.Println(1, "\tNo response for Status of %d & Content-Type of \"%s\"", status, CT)
}

// ResponseHasWrongStatus --
func (log Nice) ResponseHasWrongStatus(resp *api.Response, actualStatus int) {
	log.Println(1, "\tExpected the %d status in response, but got %d.", resp.StatusCode, actualStatus)
}

// ResponseHasWrongContentType --
func (log Nice) ResponseHasWrongContentType(resp *api.Response, actualCT string) {
	log.Println(1, "\tExpected the \"%s\" Content-Type in response, but got \"%s\".", resp.ContentType, actualCT)
}

// UsingRequest --
/* func (log Nice) UsingRequest(req *api.Request) {
	log.Println(1, "\tUsing the \"%s\" request.", req.ContentType)
} */

// UsingResponse --
func (log Nice) UsingResponse(resp *api.Response) {
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
func (log Nice) HeaderHasNoValue(hdr *api.Header) {
	log.Println(1, "\tHeader \"%s\" is required but is not present.", hdr.Name)
}

// HeaderHasWrongType --
/* func (log Nice) HeaderHasWrongType(hdr *api.Header) {
	log.Println(1, "\tHeader \"%s\" has a wrong type.", hdr.Name)
} */

// OperationOK --
func (log Nice) OperationOK(res *api.Operation) {
	log.Print(2, "\t")
	log.Println(1, "%s", log.colorOK("OK"))
}

// OperationFail --
func (log Nice) OperationFail(res *api.Operation) {
	log.Print(2, "\t")
	log.Println(1, "%s", log.colorFailure("FAILURE"))
}

// OperationNotFound --
func (log Nice) OperationNotFound(op string) {
	log.Println(1, "The operation \"%s\" isn't there.", op)
}

// SchemaTesting --
func (log Nice) SchemaTesting(schema *api.Schema, data interface{}) {
	datas := log.colorValue(fmt.Sprintf("%#v", data))
	log.Print(4, "\t%s: testing %s", log.colorID(schema.Name), datas)
}

// SchemaOK --
func (log Nice) SchemaOK(schema *api.Schema) {
	log.Println(4, log.colorSuccess(" - OK"))
}

// SchemaFail --
func (log Nice) SchemaFail(schema *api.Schema, errors []gojsonschema.ResultError) {
	log.Println(4, log.colorError(" - FAILURE"))
	// log.Println(4, "\tSchema \"%s\" has errors.", schema.Name)

	for _, desc := range errors {
		log.Println(4, "\t\t%s", log.colorError(desc))
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
func (log Nice) ParameterHasNoExample(paramName string, in string, container string) {
	log.Println(5, "\tThe %s parameter \"%s\" (from %s) has no example value to use.", in, paramName, container)
}

// UsingParameterExample --
func (log Nice) UsingParameterExample(paramName string, in string, container string) {
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
func (log Nice) TestingProject(pi *api.ProjectInfo) {
	log.Println(1, "Testing the %s @ %s", pi.Title, pi.Version)
}

// UsingHost --
func (log Nice) UsingHost(host *api.Host) {
	log.Println(2, "Using the \"%s\" host @ %s", host.Name, host.URL)
}

// TestingOperation --
func (log Nice) TestingOperation(op *api.Operation) {
	log.Print(1, "Testing the \"%s\" operation... ", log.colorOp(op.Name))
	log.Print(2, "\n")
}
