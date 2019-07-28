package test

import (
	"github.com/x1n13y84issmd42/goasis/src/api"
	"github.com/x1n13y84issmd42/goasis/src/log"
	"github.com/x1n13y84issmd42/goasis/src/utility"
)

// SchemaNumber asserts that the provided data is a number.
type SchemaNumber struct {
	Schema *api.Schema
	Log    log.ILogger
}

// Test tests.
func (schema SchemaNumber) Test(v interface{}, ctx *utility.Context) (isit bool) {
	_, isit = v.(int64)
	// fmt.Printf("SchemaNumber int64 result: %v\n", isit)

	if !isit {
		_, isit = v.(float64)
		// fmt.Printf("SchemaNumber float64 result: %v\n", isit)
	}

	return
}
