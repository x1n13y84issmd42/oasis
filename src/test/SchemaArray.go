package test

import (
	"github.com/x1n13y84issmd42/goasis/src/api"
	"github.com/x1n13y84issmd42/goasis/src/log"
	"github.com/x1n13y84issmd42/goasis/src/utility"
)

// SchemaArray asserts that the provided data is a string.
type SchemaArray struct {
	Schema *api.Schema
	Log    log.ILogger
}

// Test tests.
func (schema SchemaArray) Test(v interface{}, ctx *utility.Context) bool {
	va, isit := v.([]interface{})
	// fmt.Printf("SchemaArray result: %v\n", isit)

	test := Schema{schema.Schema.Items, schema.Log}
	OK := isit

	for vI, item := range va {
		tr := test.Test(item, ctx.PushIndex(vI, item))
		OK = OK && tr
	}

	return OK
}
