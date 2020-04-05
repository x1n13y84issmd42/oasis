package api

import "strconv"

// JSONSchema is an internal type to hold a JSON schema definition.
type JSONSchema map[string]interface{}

// Schema represents a JSON Schema used to validate response data.
type Schema struct {
	Name       string
	JSONSchema JSONSchema
}

// Cast attempts to cast a string value to a native type
// according to the provided JSON schema.
func (schema *Schema) Cast(v string) interface{} {
	switch schema.JSONSchema["type"] {
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
