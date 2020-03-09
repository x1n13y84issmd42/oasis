package test

import (
	"github.com/x1n13y84issmd42/goasis/src/api"
	"github.com/x1n13y84issmd42/goasis/src/log"
	"github.com/x1n13y84issmd42/goasis/src/utility"
)

// Schema validates structured data against API schema definitions.
type Schema struct {
	APISchema *api.Schema
	Log       log.ILogger
}

// Test doesn't do much by itself, it relies on type-specific schema tests.
func (test Schema) Test(data interface{}, ctx *utility.Context) bool {
	switch test.APISchema.DataType {
	case api.DataTypeString:
		return SchemaString{test.APISchema, test.Log}.Test(data, ctx)

	case api.DataTypeNumber:
		return SchemaNumber{test.APISchema, test.Log}.Test(data, ctx)

	case api.DataTypeInteger:
		return SchemaInteger{test.APISchema, test.Log}.Test(data, ctx)

	case api.DataTypeBoolean:
		return SchemaBoolean{test.APISchema, test.Log}.Test(data, ctx)

	case api.DataTypeArray:
		return SchemaArray{test.APISchema, test.Log}.Test(data, ctx)

	case api.DataTypeObject:
		return SchemaObject{test.APISchema, test.Log}.Test(data, ctx)

	default:
		test.Log.UnknownSchemaDataType(test.APISchema)
	}
	return false
}
