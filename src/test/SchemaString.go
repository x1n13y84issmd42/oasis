package test

import (
	"github.com/x1n13y84issmd42/goasis/src/api"
	"github.com/x1n13y84issmd42/goasis/src/log"
	"github.com/x1n13y84issmd42/goasis/src/utility"
)

// SchemaString asserts that the provided data is a string.
type SchemaString struct {
	Schema *api.Schema
	Log    log.ILogger
}

// Test tests.
func (schema SchemaString) Test(v interface{}, ctx *utility.Context) (isit bool) {
	_, isit = v.(string)
	// fmt.Printf("SchemaString result: %v\n", isit)
	return
}
