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
	// fmt.Printf("\tBetter luck next time.\n")
	fmt.Printf("\tError: %s\n", err.Error())
}

// Overriding --
func (log Simple) Overriding(what string) {
	fmt.Printf("\tOverriding %s.\n", what)
}

// ResponseHasWrongStatus --
func (log Simple) ResponseHasWrongStatus(resp *api.Response, actualStatus int) {
	fmt.Printf("\tExpected the %d status in response, but got %d.\n", resp.StatusCode, actualStatus)
}

// ResponseHasWrongContentType --
func (log Simple) ResponseHasWrongContentType(resp *api.Response, actualCT string) {
	fmt.Printf("\tExpected the '%s' Content-Type in response, but got '%s'.\n", resp.ContentType, actualCT)
}

// UsingRequest --
func (log Simple) UsingRequest(req *api.Request) {
	fmt.Printf("\tUsing the '%s' request.\n", req.ContentType)
}

// UsingResponse --
func (log Simple) UsingResponse(resp *api.Response) {
	if resp.Schema != nil {
		fmt.Printf("\tTesting against the '%s' response.\n", resp.Schema.Name)
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
	fmt.Printf("\tHeader '%s' is required but is not present.\n", hdr.Name)
}

// HeaderHasWrongType --
func (log Simple) HeaderHasWrongType(hdr *api.Header) {
	fmt.Printf("\tHeader '%s' has a wrong type.\n", hdr.Name)
}

// OperationOK --
func (log Simple) OperationOK(res *api.Operation) {
	fmt.Printf("OK\n")
}

// OperationFail --
func (log Simple) OperationFail(res *api.Operation) {
	fmt.Printf("FAILURE\n")
}

// SchemaOK --
func (log Simple) SchemaOK(schema *api.Schema) {
}

// SchemaFail --
func (log Simple) SchemaFail(schema *api.Schema) {
	fmt.Printf("\tSchema '%s' has errors.\n", schema.Name)
}

// UsingSecurity --
func (log Simple) UsingSecurity(sec *api.Security) {
	fmt.Printf("\tUsing the '%s' security settings.\n", sec.Name)
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
	fmt.Printf("Testing the \"%s\" operation @ %s %s\n", res.Name, strings.ToUpper(res.Method), res.Path)
}
