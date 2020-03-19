package test

import (
	"strconv"

	"github.com/x1n13y84issmd42/oasis/src/api"
	"github.com/x1n13y84issmd42/oasis/src/log"
	"github.com/x1n13y84issmd42/oasis/src/utility"
)

// SchemaInteger asserts that the provided data is an integer number.
type SchemaInteger struct {
	APISchema *api.Schema
	Log       log.ILogger
}

// Test checks if the passed value v is an integer one.
// Two cases are possible:
//		1. a string arrives (usually from a response header)
//			try to ParseInt() it
//		2. a native value arrives (usually from a parsed JSON)
//			just type cast it and see what happens
func (test SchemaInteger) Test(v interface{}, ctx *utility.Context) (isit bool) {
	sv, isstring := v.(string)
	if isstring {
		_, err := strconv.ParseInt(sv, 10, 64)
		if err == nil {
			isit = true
		}
	} else {
		_, isit = v.(int64)
	}

	if !isit {
		test.Log.SchemaExpectedInteger(test.APISchema, v)
	}

	return
}
