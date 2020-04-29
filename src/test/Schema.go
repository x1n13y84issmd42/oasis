package test

import (
	"github.com/x1n13y84issmd42/oasis/src/api"
	"github.com/x1n13y84issmd42/oasis/src/log"
	"github.com/xeipuuv/gojsonschema"
)

// Schema validates data against the provided schema.
func Schema(data interface{}, schema *api.Schema, logger log.ILogger) bool {
	logger.SchemaTesting(schema, data)

	schemaLoader := gojsonschema.NewGoLoader(schema.JSONSchema)
	dataLoader := gojsonschema.NewGoLoader(data)

	result, err := gojsonschema.Validate(schemaLoader, dataLoader)
	if err != nil {
		logger.Error(err)
		return false
	}

	if result.Valid() {
		logger.SchemaOK(schema)
		return true
	}

	logger.SchemaFail(schema, result.Errors())

	return false
}
