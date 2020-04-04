package utility

import (
	"strconv"

	"github.com/x1n13y84issmd42/oasis/src/api"
)

// SchemaCast attempts to cast a string value to a native type
// according to the provided JSON schema.
func SchemaCast(v string, sch *api.Schema) interface{} {
	switch sch.JSONSchema["type"] {
	case "number":
		r, err := strconv.ParseFloat(v, 64)
		if err == nil {
			return r
		}

	case "integer":
		r, err := strconv.ParseInt(v, 10, 64)
		if err == nil {
			return r
		}

	case "boolean":
		r, err := strconv.ParseBool(v)
		if err == nil {
			return r
		}
	}

	return v
}
