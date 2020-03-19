package test

import (
	"github.com/x1n13y84issmd42/oasis/src/api"
	"github.com/x1n13y84issmd42/oasis/src/log"
	"github.com/x1n13y84issmd42/oasis/src/utility"
)

// SchemaArray asserts that the provided data is an array.
type SchemaArray struct {
	APISchema *api.Schema
	Log       log.ILogger
}

// Test ensures that v is an array, iterates over items and validate them against them schema.
func (test SchemaArray) Test(v interface{}, ctx *utility.Context) bool {
	items, isit := v.([]interface{})
	OK := isit

	if !OK {
		test.Log.SchemaExpectedArray(test.APISchema, v)
		return false
	}

	testItems := Schema{test.APISchema.Items, test.Log}

	for vI, item := range items {
		OK = OK && testItems.Test(item, ctx.PushIndex(vI, item))
	}

	return OK
}
