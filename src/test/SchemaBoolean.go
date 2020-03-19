package test

import (
	"strconv"

	"github.com/x1n13y84issmd42/oasis/src/api"
	"github.com/x1n13y84issmd42/oasis/src/log"
	"github.com/x1n13y84issmd42/oasis/src/utility"
)

// SchemaBoolean asserts that the provided data is a boolean value.
type SchemaBoolean struct {
	APISchema *api.Schema
	Log       log.ILogger
}

// Test checks if the passed value v is a boolean one.
// Two cases are possible:
//		1. a string arrives (usually from a response header)
//			try to ParseBool() it
//		2. a native value arrives (usually from a parsed JSON)
//			just type cast it and see what happens
func (test SchemaBoolean) Test(v interface{}, ctx *utility.Context) (isit bool) {
	sv, isstring := v.(string)
	if isstring {
		_, err := strconv.ParseBool(sv)
		if err == nil {
			isit = true
		}
	} else {
		_, isit = v.(bool)
	}

	if !isit {
		test.Log.SchemaExpectedBoolean(test.APISchema, v)
	}

	return
}
