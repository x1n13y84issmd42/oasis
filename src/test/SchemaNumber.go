package test

import (
	"strconv"

	"github.com/x1n13y84issmd42/oasis/src/api"
	"github.com/x1n13y84issmd42/oasis/src/log"
	"github.com/x1n13y84issmd42/oasis/src/utility"
)

// SchemaNumber asserts that the provided data is a floating point number.
type SchemaNumber struct {
	APISchema *api.Schema
	Log       log.ILogger
}

// Test checks if the passed value v is a floating point one.
// Two cases are possible:
//		1. a string arrives (usually from a response header)
//			try to ParseFloat() it
//		2. a native value arrives (usually from a parsed JSON)
//			just type cast it and see what happens
func (test SchemaNumber) Test(v interface{}, ctx *utility.Context) (isNumber bool) {
	sv, isstring := v.(string)
	if isstring {
		_, err := strconv.ParseFloat(sv, 64)
		if err == nil {
			isNumber = true
		}
	} else {
		_, isNumber = v.(float64)
	}

	if !isNumber {
		test.Log.SchemaExpectedNumber(test.APISchema, v)
	}

	return
}
