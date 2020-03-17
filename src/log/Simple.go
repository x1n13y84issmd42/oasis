package log

import (
	"fmt"
	"strings"

	"github.com/x1n13y84issmd42/goasis/src/api"
	"github.com/x1n13y84issmd42/goasis/src/utility"
)

// Simple - a simple test execution logger
type Simple struct{}

// UsingDefaultHost --
func (log Simple) UsingDefaultHost() {
	fmt.Printf("No host name has been specified, using the first one in the list.\n")
}

func (log Simple) HostNotFound(h string) {
	if h == "" {
		fmt.Printf("No default host is found in the spec.\n")
	} else {
		fmt.Printf("The host \"%s\" is not found in the spec.\n", h)
	}
}

// ResponseExpectedArray --
func (log Simple) ResponseExpectedArray(schema *api.Response) {
	fmt.Printf("\tExpected an array in response, but received something else.\n")
}

// ResponseExpectedObject --
func (log Simple) ResponseExpectedObject(schema *api.Response) {
	fmt.Printf("\tExpected an object in response, but received something else.\n")
}

// Error --
func (log Simple) Error(err error) {
	fmt.Printf("\tSomething happened: %s\n", err.Error())
}

// Overriding --
func (log Simple) Overriding(what string) {
	fmt.Printf("\tOverriding %s.\n", what)
}

// Requesting --
func (log Simple) Requesting(URL string) {
	fmt.Printf("\tRequesting %s\n", URL)
}

// ResponseNotFound --
func (log Simple) ResponseNotFound(CT string, status int) {
	fmt.Printf("\tNo response for Status of %d & Content-Type of \"%s\"\n", status, CT)
}

// ResponseHasWrongStatus --
func (log Simple) ResponseHasWrongStatus(resp *api.Response, actualStatus int) {
	fmt.Printf("\tExpected the %d status in response, but got %d.\n", resp.StatusCode, actualStatus)
}

// ResponseHasWrongContentType --
func (log Simple) ResponseHasWrongContentType(resp *api.Response, actualCT string) {
	fmt.Printf("\tExpected the \"%s\" Content-Type in response, but got \"%s\".\n", resp.ContentType, actualCT)
}

// UsingRequest --
func (log Simple) UsingRequest(req *api.Request) {
	fmt.Printf("\tUsing the \"%s\" request.\n", req.ContentType)
}

// UsingResponse --
func (log Simple) UsingResponse(resp *api.Response) {
	if resp.Schema != nil {
		fmt.Printf("\tTesting against the \"%s\" response.\n", resp.Schema.Name)
	} else {
		CT := resp.ContentType
		if len(CT) == 0 {
			CT = "*/*"
		}
		fmt.Printf("\tTesting against the %s @ %d response.\n", CT, resp.StatusCode)
	}
}

// HeaderHasNoValue --
func (log Simple) HeaderHasNoValue(hdr *api.Header) {
	fmt.Printf("\tHeader \"%s\" is required but is not present.\n", hdr.Name)
}

// HeaderHasWrongType --
func (log Simple) HeaderHasWrongType(hdr *api.Header) {
	fmt.Printf("\tHeader \"%s\" has a wrong type.\n", hdr.Name)
}

// OperationOK --
func (log Simple) OperationOK(res *api.Operation) {
	fmt.Printf("OK\n")
}

// OperationFail --
func (log Simple) OperationFail(res *api.Operation) {
	fmt.Printf("FAILURE\n")
}

// OperationNotFound --
func (log Simple) OperationNotFound(op string) {
	fmt.Printf("The operation \"%s\" isn't there.\n", op)
}

// SchemaOK --
func (log Simple) SchemaOK(schema *api.Schema) {
}

// SchemaFail --
func (log Simple) SchemaFail(schema *api.Schema) {
	fmt.Printf("\tSchema \"%s\" has errors.\n", schema.Name)
}

// UnknownSchemaDataType --
func (log Simple) UnknownSchemaDataType(schema *api.Schema) {
	fmt.Printf("\tSchema \"%s\" has unknown data type \"%s\".\n", schema.Name, schema.DataType)
}

// SchemaExpectedBoolean --
func (log Simple) SchemaExpectedBoolean(schema *api.Schema, v interface{}) {
	fmt.Printf("\tSchema \"%s\" expected %#v to be a boolean type.\n", schema.Name, v)
}

// SchemaExpectedNumber --
func (log Simple) SchemaExpectedNumber(schema *api.Schema, v interface{}) {
	fmt.Printf("\tSchema \"%s\" expected %#v to be a floating point number.\n", schema.Name, v)
}

// SchemaExpectedInteger --
func (log Simple) SchemaExpectedInteger(schema *api.Schema, v interface{}) {
	fmt.Printf("\tSchema \"%s\" expected %#v to be an integer number.\n", schema.Name, v)
}

// SchemaExpectedString --
func (log Simple) SchemaExpectedString(schema *api.Schema, v interface{}) {
	fmt.Printf("\tSchema \"%s\" expected %#v to be a string type.\n", schema.Name, v)
}

// SchemaExpectedArray --
func (log Simple) SchemaExpectedArray(schema *api.Schema, v interface{}) {
	fmt.Printf("\tSchema \"%s\" expected %#v to be an array type.\n", schema.Name, v)
}

// SchemaExpectedObject --
func (log Simple) SchemaExpectedObject(schema *api.Schema, v interface{}) {
	fmt.Printf("\tSchema \"%s\" expected %#v to be an object type.\n", schema.Name, v)
}

// UsingSecurity --
func (log Simple) UsingSecurity(sec *api.Security) {
	fmt.Printf("\tUsing the \"%s\" security settings.\n", sec.Name)
}

// ParameterHasNoExample --
func (log Simple) ParameterHasNoExample(param *api.Parameter, container string) {
	fmt.Printf("\tThe %s parameter \"%s\" has no example value to use.\n", container, param.Name)
}

// UsingParameterExample --
func (log Simple) UsingParameterExample(param *api.Parameter, container string) {
	fmt.Printf("\tUsing the %s parameter \"%s\" example.\n", container, param.Name)
}

// PropertyHasNoValue --
func (log Simple) PropertyHasNoValue(prop *api.Property, ctx *utility.Context) {
	fmt.Printf("\t%s: property is required but is not present.\n", ctx.String())
}

// PropertyHasWrongType --
func (log Simple) PropertyHasWrongType(prop *api.Property, ctx *utility.Context) {
	fmt.Printf("\t%s: property has wrong type. Expected %s, got %s.\n", ctx.String(), prop.Schema.DataType, ctx.CurrentValueType())
}

// TestingProject --
func (log Simple) TestingProject(pi *api.ProjectInfo) {
	fmt.Printf("Testing the %s @ %s\n", pi.Title, pi.Version)
}

// UsingHost --
func (log Simple) UsingHost(host *api.Host) {
	fmt.Printf("Using the \"%s\" host @ %s\n", host.Name, host.URL)
}

// TestingOperation --
func (log Simple) TestingOperation(res *api.Operation) {
	fmt.Printf("Testing the \"%s\" operation @ %s %s\n", res.Name, strings.ToUpper(res.Method), res.Path.Path)
}
