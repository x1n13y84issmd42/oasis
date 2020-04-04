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
		fmt.Printf("The document is valid\n")
		return true
	} else {
		fmt.Printf("The document is not valid. see errors :\n")
		for _, desc := range result.Errors() {
			fmt.Printf("- %s\n", desc)
		}
	}

	return false
}
