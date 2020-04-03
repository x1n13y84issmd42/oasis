package api

// JSONSchema is an internal type to hold a JSON schema definition.
type JSONSchema map[string]interface{}

// Schema represents a JSON Schema used to validate response data.
type Schema struct {
	Name       string
	JSONSchema JSONSchema
}
