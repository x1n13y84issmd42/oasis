package test

import (
	"fmt"

	"github.com/x1n13y84issmd42/goasis/src/api"
	"github.com/x1n13y84issmd42/goasis/src/log"
	"github.com/x1n13y84issmd42/goasis/src/utility"
)

// Schema validates structured data against API schema definitions.
type Schema struct {
	APISchema *api.Schema
	Log       log.ILogger
}

// Test tests.
func (schema Schema) Test(data interface{}, ctx *utility.Context) bool {
	switch schema.APISchema.DataType {
	case api.DataTypeString:
		return SchemaString{schema.APISchema, schema.Log}.Test(data, ctx)

	case api.DataTypeNumber:
		return SchemaNumber{schema.APISchema, schema.Log}.Test(data, ctx)

	case api.DataTypeBoolean:
		return SchemaBoolean{schema.APISchema, schema.Log}.Test(data, ctx)

	case api.DataTypeArray:
		return SchemaArray{schema.APISchema, schema.Log}.Test(data, ctx)

	case api.DataTypeObject:
		return SchemaObject{schema.APISchema, schema.Log}.Test(data, ctx)

	default:
		fmt.Println("Woot")
		schema.Log.UnknownSchemaDataType(schema.APISchema)
	}
	return false
}
