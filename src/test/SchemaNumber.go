package test

import (
	"fmt"
	"strconv"

	"github.com/x1n13y84issmd42/goasis/src/api"
	"github.com/x1n13y84issmd42/goasis/src/log"
	"github.com/x1n13y84issmd42/goasis/src/utility"
)

// SchemaNumber asserts that the provided data is a number.
type SchemaNumber struct {
	APISchema *api.Schema
	Log       log.ILogger
}

// Test checks if the passed value v is a numeric one.
// Two cases are possible:
//		1. a string arrives (usually from a response header)
//			try to ParseInt() it, if it fails - try ParseFloat()
//		2. a native value arrives (usually from a parsed JSON)
//			just type cast it and see what happens
func (test SchemaNumber) Test(v interface{}, ctx *utility.Context) (isit bool) {
	sv, isstring := v.(string)
	if isstring {
		_, err := strconv.ParseInt(sv, 10, 64)
		if err != nil {
			_, err := strconv.ParseFloat(sv, 64)
			if err == nil {
				isit = true
			}
		} else {
			isit = true
		}
	} else {
		_, isit = v.(int64)
		fmt.Printf("SchemaNumber int64 result for \"%s\": %v\n", v.(string), isit)

		if !isit {
			_, isit = v.(float64)
			fmt.Printf("SchemaNumber float64 result \"%s\": %v\n", v.(string), isit)
		}
	}

	if !isit {
		test.Log.SchemaExpectedNumber(test.APISchema, v)
	}

	return
}
