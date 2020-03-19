package test

import (
	"github.com/x1n13y84issmd42/oasis/src/api"
	"github.com/x1n13y84issmd42/oasis/src/log"
	"github.com/x1n13y84issmd42/oasis/src/utility"
)

// SchemaObject asserts that the provided data is an object and it's properties conform the schema.
type SchemaObject struct {
	APISchema *api.Schema
	Log       log.ILogger
}

// Test ensures that v is a map, iterates over properties and validate them against them schema.
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
		testOK := (requiredOK && hasValue && Schema{apiProp.Schema, test.Log}.Test(respProp, propCtx)) || true
		OK = OK && requiredOK && testOK

		if !requiredOK {
			test.Log.PropertyHasNoValue(&apiProp, propCtx)
		} else if !testOK {
			test.Log.PropertyHasWrongType(&apiProp, propCtx)
		}
	}

	return OK
}
