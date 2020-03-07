package test

import (
	"github.com/x1n13y84issmd42/goasis/src/api"
	"github.com/x1n13y84issmd42/goasis/src/log"
	"github.com/x1n13y84issmd42/goasis/src/utility"
)

// SchemaObject asserts that the provided data is an object.
type SchemaObject struct {
	APISchema *api.Schema
	Log       log.ILogger
}

// Test tests.
func (test SchemaObject) Test(v interface{}, ctx *utility.Context) bool {
	resp, isit := v.(map[string]interface{})

	OK := isit

	if !OK {
		test.Log.SchemaExpectedObject(test.APISchema, v)
		return false
	}

	for _, apiProp := range *test.APISchema.Properties {
		respProp := resp[apiProp.Name]
		propCtx := ctx.PushProperty(apiProp.Name, respProp)
		hasValue := respProp != nil
		requiredOK := (apiProp.Required && hasValue) || !apiProp.Required
		testOK := Schema{apiProp.Schema, test.Log}.Test(respProp, propCtx)
		OK = OK && requiredOK && testOK

		if !requiredOK {
			test.Log.PropertyHasNoValue(&apiProp, propCtx)
		} else if !testOK {
			test.Log.PropertyHasWrongType(&apiProp, propCtx)
		}
	}

	return OK
}
