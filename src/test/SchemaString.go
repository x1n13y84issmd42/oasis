package test

import (
	"github.com/x1n13y84issmd42/goasis/src/api"
	"github.com/x1n13y84issmd42/goasis/src/log"
	"github.com/x1n13y84issmd42/goasis/src/utility"
)

// SchemaString asserts that the provided data is a string.
type SchemaString struct {
	APISchema *api.Schema
	Log       log.ILogger
}

// Test tests.
func (test SchemaString) Test(v interface{}, ctx *utility.Context) (isit bool) {
	_, isit = v.(string)

	if !isit {
		test.Log.SchemaExpectedString(test.APISchema, v)
	}

	return
}
