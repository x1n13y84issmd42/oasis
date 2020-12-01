package params

import (
	goerrors "errors"
	"regexp"
	"strconv"

	"github.com/x1n13y84issmd42/oasis/src/contract"
	"github.com/x1n13y84issmd42/oasis/src/errors"
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
			return pr.Cast(access(res, pr.Log))
		}

		if res, err := test.TryJSONArrayResponse(&pr.Result.ResponseBytes, pr.Log); err == nil {
			return pr.Cast(access(res, pr.Log))
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

		return pr.Cast(access(data, pr.Log))
	}
}

// Cast casts the given value to string.
func (pr Reference) Cast(v interface{}) string {
	return Cast(v)
}

// ReferenceAccess is a function to compute and return a referenced value.
type ReferenceAccess func(interface{}, contract.Logger) interface{}

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

// ParseSelector parses selectors and constructs a parameter access function from it.
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
			return NoAccess(goerrors.New("impossible to parse selector " + selector)), selector
		}
	}

	return res, ""
}

// AccessContent passes through.
func AccessContent() ReferenceAccess {
	return func(v interface{}, log contract.Logger) interface{} {
		return v
	}
}

// AccessArray treats v as an array and returns i-th element of it.
func AccessArray(access ReferenceAccess, i int) ReferenceAccess {
	return func(v interface{}, log contract.Logger) interface{} {
		xv := access(v, log)
		arrv, ok := xv.([]interface{})
		if ok {
			if len(arrv) <= i {
				errors.Report(errors.OutOfRange(i, &arrv, nil), "AccessArray", log)
			}

			return (arrv)[i]
		}

		errors.Report(errors.NotAn("array", &arrv, nil), "AccessArray", log)
		return nil
	}
}

// AccessObject treats v as a map and returns f-th element of it.
func AccessObject(access ReferenceAccess, f string) ReferenceAccess {
	return func(v interface{}, log contract.Logger) interface{} {
		objv, ok := access(v, log).(map[string]interface{})
		if ok {
			v, vok := (objv)[f]
			if vok {
				return v
			}

			errors.Report(errors.NoProperty(f, nil), "AccessObject", log)
		}

		errors.Report(errors.NotAn("object", &objv, nil), "AccessObject", log)
		return nil
	}
}

// NoAccess is a placeholder ref access function used when we can't have a real one.
// Usually when it is impossible to parse a reference selector.
func NoAccess(err error) ReferenceAccess {
	return func(v interface{}, log contract.Logger) interface{} {
		if err != nil {
			log.Error(err)
		}

		errors.Report(err, "NoAccess", log)
		return nil
	}
}
