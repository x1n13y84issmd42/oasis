package params

import (
	"fmt"
	"strconv"
)

// Value is the default pass-through function to provide parameters.
func Value(v string) func() string {
	return func() string {
		return v
	}
}

// Cast casts the given value to string.
func Cast(v interface{}) string {
	if cv, ok := v.(string); ok {
		return cv
	}

	if cv, ok := v.(int64); ok {
		return strconv.Itoa(int(cv))
	}

	if cv, ok := v.(float64); ok {
		// This is here because json.Unmarshal parses integer values as float64s.
		if float64(int64(cv)) == cv {
			return fmt.Sprintf("%d", int64(cv))
		}
		//TODO: this loses precision on numbers like .00000001
		res := fmt.Sprintf("%f", cv)
		return res
	}

	if cv, ok := v.(bool); ok {
		if cv {
			return "true"
		}

		return "false"
	}

	//TODO: this is very questionable :/
	return fmt.Sprintf("%#v", v)
}
