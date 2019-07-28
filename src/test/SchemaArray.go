package test

import (
	"github.com/x1n13y84issmd42/goasis/src/api"
	"github.com/x1n13y84issmd42/goasis/src/log"
	"github.com/x1n13y84issmd42/goasis/src/utility"
)

// SchemaArray asserts that the provided data is an array.
type SchemaArray struct {
	Schema *api.Schema
	Log    log.ILogger
}

// Test tests.
func (schema SchemaArray) Test(v interface{}, ctx *utility.Context) bool {
	items, isit := v.([]interface{})
	OK := isit

	if !OK {
		return false
	}

	test := Schema{schema.Schema.Items, schema.Log}

	for vI, item := range items {
		OK = OK && test.Test(item, ctx.PushIndex(vI, item))
	}

	return OK
}
