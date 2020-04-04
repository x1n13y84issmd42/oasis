package test

import (
	"fmt"

	"github.com/x1n13y84issmd42/oasis/src/api"
	"github.com/xeipuuv/gojsonschema"
)

// ValidateJSONSchema uses the provided schema to validate the provided data against.
func ValidateJSONSchema(data interface{}, schema *api.Schema) bool {
	fmt.Printf("\tTesting %#v\n", data)

	schemaLoader := gojsonschema.NewGoLoader(schema.JSONSchema)
	dataLoader := gojsonschema.NewGoLoader(data)

	result, err := gojsonschema.Validate(schemaLoader, dataLoader)
	if err != nil {
		panic(err.Error())
	}

	if result.Valid() {
		fmt.Printf("\tOK\n")
		return true
	}

	fmt.Printf("\tFAILURE:\n")
	for _, desc := range result.Errors() {
		fmt.Printf("\t - %s\n", desc)
	}

	return false
}
