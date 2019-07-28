package test

import (
	"github.com/x1n13y84issmd42/goasis/src/api"
	"github.com/x1n13y84issmd42/goasis/src/log"
	"github.com/x1n13y84issmd42/goasis/src/utility"
)

// SchemaBoolean asserts that the provided data is a string.
type SchemaBoolean struct {
	Schema *api.Schema
	Log    log.ILogger
}

// Test tests.
func (schema SchemaBoolean) Test(v interface{}, ctx *utility.Context) (isit bool) {
	_, isit = v.(bool)
	// fmt.Printf("SchemaBoolean result: %v\n", isit)
	return
}
