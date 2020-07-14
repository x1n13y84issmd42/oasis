package params

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/x1n13y84issmd42/oasis/src/contract"
	"github.com/x1n13y84issmd42/oasis/src/strings"
	"github.com/x1n13y84issmd42/oasis/src/test"
)

// Reference is a special kind of parameter which comes from
// an operation response. When a parameter for some operation in a script
// has a value like "operationID.response.[0].user.id" this means that
// the actual value comes from JSON response of the operation "operationID"
// and it's exact location is "[0].user.id" field.
type Reference struct {
	OpID     string
	Result   *contract.OperationResult
	Selector string

	Log contract.Logger
}

// Value returns a parameter access function which computes and returns a real value.
func (pr Reference) Value() contract.ParameterAccess {
	return func() string {
		access, _ := ParseSelector(pr.Selector, pr.Log)

		var data interface{}
		var err error

		if res, err := test.TryJSONObjectResponse(&pr.Result.ResponseBytes, pr.Log); err == nil {
			return pr.Cast(access(res))
		}

		if res, err := test.TryJSONArrayResponse(&pr.Result.ResponseBytes, pr.Log); err == nil {
			return pr.Cast(access(res))
		}

		if res, err := test.TryJSONStringResponse(&pr.Result.ResponseBytes, pr.Log); err == nil {
			data = res
		}

		if res, err := test.TryJSONNumberResponse(&pr.Result.ResponseBytes, pr.Log); err == nil {
			data = res
		}

		if res, err := test.TryJSONBooleanResponse(&pr.Result.ResponseBytes, pr.Log); err == nil {
			data = res
		}

		if err != nil {
			pr.Log.Error(err)
		}

		return pr.Cast(access(data))
	}
}

// Cast casts the given value to string.
func (pr Reference) Cast(v interface{}) string {
	if iv, ok := v.(int64); ok {
		return strconv.Itoa(int(iv))
	}

	if fv, ok := v.(float64); ok {
		if float64(int64(fv)) == fv {
			return fmt.Sprintf("%d", int64(fv))
		}
		return fmt.Sprintf("%f", fv)
	}

	if bv, ok := v.(bool); ok {
		if bv {
			return "true"
		}

		return "false"
	}

	//TODO: this is very questionable :/
	return fmt.Sprintf("%#v", v)
}

// ReferenceAccess ...
type ReferenceAccess func(interface{}) interface{}

// ParseArrayIndexRef parses the JSON array index signature [N]
// in the beginning of selector. It returns the index and number
// of parsed characters when successful, and (-1, 0) when it was not able to parse.
func ParseArrayIndexRef(selector string) (int, int) {
	rx := regexp.MustCompile("^\\[(?P<index>\\d+)\\]")
	matches := strings.RxMatches(selector, rx)
	if matches["index"] != "" {
		// Ignoring the error because regexp requires an integer.
		idx, _ := strconv.Atoi(matches["index"])
		return idx, len(matches["index"]) + 2
	}

	return -1, 0
}

// ParseObjectFieldRef parses the JSON object field signature .fieldName
// in the beginning of selector. It returns the field name and number
// of parsed characters when successful, and ("", 0) when it was not able to parse.
func ParseObjectFieldRef(selector string) (string, int) {
	rx := regexp.MustCompile("^\\.(?P<field>[a-zA-Z_]+)")
	matches := strings.RxMatches(selector, rx)
	if matches["field"] != "" {
		return matches["field"], len(matches["field"]) + 1
	}

	return "", 0
}

// AccessContent passes through.
func AccessContent() ReferenceAccess {
	return func(v interface{}) interface{} {
		return v
	}
}

// AccessArray treats v as an array and returns i-th element of it.
func AccessArray(access ReferenceAccess, i int) ReferenceAccess {
	return func(v interface{}) interface{} {
		xv := access(v)
		arrv, ok := xv.(*[]interface{})
		if ok {
			if len(*arrv) <= i {
				panic(fmt.Sprintf("Array index %d is out of range 0-%d.", i, (len(*arrv) - 1)))
			}

			return (*arrv)[i]
		}
		panic("Not an array.")
	}
}

// AccessObject treats v as a map and returns f-th element of it.
func AccessObject(access ReferenceAccess, f string) ReferenceAccess {
	return func(v interface{}) interface{} {
		objv, ok := access(v).(map[string]interface{})
		if ok {
			v, vok := objv[f]
			if vok {
				return v
			}

			panic("Field '" + f + "' is not found in the object.")
		}

		panic("Not an object.")
	}
}

// NoAccess ...
func NoAccess(err error, log contract.Logger) ReferenceAccess {
	return func(interface{}) interface{} {
		if err != nil {
			log.Error(err)
		} else {
			log.Error(errors.New("no error in a null object"))
		}
		os.Exit(1)
		return nil
	}
}

// ParseSelector parses selectors and constructs a .
func ParseSelector(selector string, log contract.Logger) (ReferenceAccess, string) {
	res := AccessContent()

	for len(selector) > 0 {
		if i, n := ParseArrayIndexRef(selector); n > 0 {
			res = AccessArray(res, i)
			selector = selector[n:]
		} else if f, n := ParseObjectFieldRef(selector); n > 0 {
			res = AccessObject(res, f)
			selector = selector[n:]
		} else {
			return NoAccess(errors.New("impossible to parse selector "+selector), log), selector
		}
	}

	return res, ""
}
