package test

import (
	"github.com/x1n13y84issmd42/oasis/src/api"
	"github.com/x1n13y84issmd42/oasis/src/log"
	"github.com/x1n13y84issmd42/oasis/src/utility"
)

// SchemaString asserts that the provided data is a string.
type SchemaString struct {
	APISchema *api.Schema
	Log       log.ILogger
}

// Test tests.
func (test SchemaString) Test(v interface{}, ctx *utility.Context) (isString bool) {
	_, isString = v.(string)

	if !isString {
		test.Log.SchemaExpectedString(test.APISchema, v)
	}

	return
}
