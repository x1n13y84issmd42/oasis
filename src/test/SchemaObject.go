package test

import (
	"github.com/x1n13y84issmd42/goasis/src/api"
	"github.com/x1n13y84issmd42/goasis/src/log"
	"github.com/x1n13y84issmd42/goasis/src/utility"
)

// SchemaObject asserts that the provided data is a string.
type SchemaObject struct {
	Schema *api.Schema
	Log    log.ILogger
}

// Test tests.
func (schema SchemaObject) Test(v interface{}, ctx *utility.Context) bool {
	resp, isit := v.(map[string]interface{})

	OK := isit

	if !OK {
		return false
	}

	for _, apiProp := range *schema.Schema.Properties {
		respProp := resp[apiProp.Name]
		propCtx := ctx.PushProperty(apiProp.Name, respProp)
		hasValue := respProp != nil
		requiredOK := (apiProp.Required && hasValue) || !apiProp.Required
		testOK := Schema{apiProp.Schema, schema.Log}.Test(respProp, propCtx)
		OK = OK && requiredOK && testOK

		if !requiredOK {
			schema.Log.PropertyHasNoValue(&apiProp, propCtx)
		} else if !testOK {
			schema.Log.PropertyHasWrongType(&apiProp, propCtx)
		}
	}

	return OK
}
