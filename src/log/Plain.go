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

// Usage prints CLI usage information.
func (log Plain) Usage() {
	fmt.Println("Please specify at least a spec file & an operation to test.")
	fmt.Println("Example:")
	fmt.Println("oasis from path/to/oas_spec.yaml test operation_id")
}

// Error --
func (log Plain) Error(err error) {
	log.Println(1, "\tSomething happened: %s", err.Error())
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

// Overriding --
func (log Plain) Overriding(what string) {
	log.Println(1, "\tOverriding %s.", what)
}

// Requesting --
func (log Plain) Requesting(method string, URL string) {
	log.Println(2, "\tRequesting %s @ %s", method, URL)
}

// ResponseNotFound --
func (log Plain) ResponseNotFound(CT string, status int) {
	log.Println(1, "\tNo response for Status of %d & Content-Type of \"%s\"", status, CT)
}

// ResponseHasWrongStatus --
func (log Plain) ResponseHasWrongStatus(resp *api.Response, actualStatus int) {
	log.Println(1, "\tExpected the %d status in response, but got %d.", resp.StatusCode, actualStatus)
}

// ResponseHasWrongContentType --
func (log Plain) ResponseHasWrongContentType(resp *api.Response, actualCT string) {
	log.Println(1, "\tExpected the \"%s\" Content-Type in response, but got \"%s\".", resp.ContentType, actualCT)
}

// UsingRequest --
/* func (log Simple) UsingRequest(req *api.Request) {
	log.Println(1, "\tUsing the \"%s\" request.", req.ContentType)
} */

// UsingResponse --
func (log Plain) UsingResponse(resp *api.Response) {
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
func (log Plain) HeaderHasNoValue(hdr *api.Header) {
	log.Println(1, "\tHeader \"%s\" is required but is not present.", hdr.Name)
}

// HeaderHasWrongType --
/* func (log Simple) HeaderHasWrongType(hdr *api.Header) {
	log.Println(1, "\tHeader \"%s\" has a wrong type.", hdr.Name)
} */

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

// UnknownSchemaDataType --
/* func (log Simple) UnknownSchemaDataType(schema *api.Schema) {
	log.Println(1, "\tSchema \"%s\" has unknown data type \"%s\".", schema.Name, schema.DataType)
} */

// SchemaExpectedBoolean --
/* func (log Simple) SchemaExpectedBoolean(schema *api.Schema, v interface{}) {
	log.Println(1, "\tSchema \"%s\" expected %#v to be a boolean type.", schema.Name, v)
} */

// SchemaExpectedNumber --
/* func (log Simple) SchemaExpectedNumber(schema *api.Schema, v interface{}) {
	log.Println(1, "\tSchema \"%s\" expected %#v to be a floating point number.", schema.Name, v)
} */

// SchemaExpectedInteger --
/* func (log Simple) SchemaExpectedInteger(schema *api.Schema, v interface{}) {
	log.Println(1, "\tSchema \"%s\" expected %#v to be an integer number.", schema.Name, v)
} */

// SchemaExpectedString --
/* func (log Simple) SchemaExpectedString(schema *api.Schema, v interface{}) {
	log.Println(1, "\tSchema \"%s\" expected %#v to be a string type.", schema.Name, v)
} */

// SchemaExpectedArray --
/* func (log Simple) SchemaExpectedArray(schema *api.Schema, v interface{}) {
	log.Println(1, "\tSchema \"%s\" expected %#v to be an array type.", schema.Name, v)
} */

// SchemaExpectedObject --
/* func (log Simple) SchemaExpectedObject(schema *api.Schema, v interface{}) {
	log.Println(1, "\tSchema \"%s\" expected %#v to be an object type.", schema.Name, v)
} */

// UsingSecurity --
func (log Plain) UsingSecurity(sec api.ISecurity) {
	log.Println(1, "\tUsing the \"%s\" security settings.", sec.GetName())
}

// ParameterHasNoExample --
func (log Plain) ParameterHasNoExample(paramName string, in string, container string) {
	log.Println(5, "\tThe %s parameter \"%s\" (from %s) has no example value to use.", in, paramName, container)
}

// UsingParameterExample --
func (log Plain) UsingParameterExample(paramName string, in string, container string) {
	log.Println(5, "\tUsing the %s parameter \"%s\" (from %s) example.", in, paramName, container)
}

// PropertyHasNoValue --
/* func (log Simple) PropertyHasNoValue(prop *api.Property, ctx *utility.Context) {
	log.Println(1, "\t%s: property is required but is not present.", ctx.String())
} */

// PropertyHasWrongType --
/* func (log Simple) PropertyHasWrongType(prop *api.Property, ctx *utility.Context) {
	log.Println(1, "\t%s: property has wrong type. Expected %s, got %s.", ctx.String(), prop.Schema.DataType, ctx.CurrentValueType())
} */

// TestingProject --
func (log Plain) TestingProject(pi *api.ProjectInfo) {
	log.Println(1, "Testing the %s @ %s", pi.Title, pi.Version)
}

// UsingHost --
func (log Plain) UsingHost(host *api.Host) {
	log.Println(2, "Using the \"%s\" host @ %s", host.Name, host.URL)
}
