package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/x1n13y84issmd42/oasis/src/api"
	"github.com/x1n13y84issmd42/oasis/src/contract"
	"github.com/x1n13y84issmd42/oasis/src/log"
)

func TestJSONResponse_String(T *testing.T) {
	result := &contract.OperationResult{
		ResponseBytes: []byte("\"anything works until it's quoted\""),
	}

	schema := &api.Schema{
		JSONSchema: api.JSONSchema{
			"type": "string",
		},
	}

	actual := JSONResponse(result, schema, log.NewFestive(0))
	assert.True(T, actual)
}

func TestJSONResponse_String_False(T *testing.T) {
	result := &contract.OperationResult{
		ResponseBytes: []byte("unquoted malformed string\""),
	}

	schema := &api.Schema{
		JSONSchema: api.JSONSchema{
			"type": "string",
		},
	}

	actual := JSONResponse(result, schema, log.NewFestive(0))
	assert.False(T, actual)
}

func TestJSONResponse_Number(T *testing.T) {
	result := &contract.OperationResult{
		ResponseBytes: []byte("42"),
	}

	schema := &api.Schema{
		JSONSchema: api.JSONSchema{
			"type": "number",
		},
	}

	actual := JSONResponse(result, schema, log.NewFestive(0))
	assert.True(T, actual)
}

func TestJSONResponse_Number_False(T *testing.T) {
	result := &contract.OperationResult{
		ResponseBytes: []byte("not a number, not even close"),
	}

	schema := &api.Schema{
		JSONSchema: api.JSONSchema{
			"type": "number",
		},
	}

	actual := JSONResponse(result, schema, log.NewFestive(0))
	assert.False(T, actual)
}

func TestJSONResponse_Boolean(T *testing.T) {
	result := &contract.OperationResult{
		ResponseBytes: []byte("false"),
	}

	schema := &api.Schema{
		JSONSchema: api.JSONSchema{
			"type": "boolean",
		},
	}

	actual := JSONResponse(result, schema, log.NewFestive(0))
	assert.True(T, actual)
}

func TestJSONResponse_Boolean_False(T *testing.T) {
	result := &contract.OperationResult{
		ResponseBytes: []byte("truth"),
	}

	schema := &api.Schema{
		JSONSchema: api.JSONSchema{
			"type": "boolean",
		},
	}

	actual := JSONResponse(result, schema, log.NewFestive(0))
	assert.False(T, actual)
}

func TestJSONResponse_Object(T *testing.T) {
	result := &contract.OperationResult{
		ResponseBytes: []byte("{\"name\":\"johnny\",\"age\":42}"),
	}

	schema := &api.Schema{
		JSONSchema: api.JSONSchema{
			"type": "object",
			"properties": map[string]api.JSONSchema{
				"name": {
					"type": "string",
				},
				"age": {
					"type": "number",
				},
			},
		},
	}

	actual := JSONResponse(result, schema, log.NewFestive(0))
	assert.True(T, actual)
}

func TestJSONResponse_Object_False_Schema(T *testing.T) {
	result := &contract.OperationResult{
		ResponseBytes: []byte("{\"name\":false,\"age\":\"42\"}"),
	}

	schema := &api.Schema{
		JSONSchema: api.JSONSchema{
			"type": "object",
			"properties": map[string]api.JSONSchema{
				"name": {
					"type": "string",
				},
				"age": {
					"type": "number",
				},
			},
		},
	}

	actual := JSONResponse(result, schema, log.NewFestive(0))
	assert.False(T, actual)
}

func TestJSONResponse_Object_False_Unmarshal(T *testing.T) {
	result := &contract.OperationResult{
		ResponseBytes: []byte("... INVALID JSON }"),
	}

	schema := &api.Schema{
		JSONSchema: api.JSONSchema{
			"type": "object",
			"properties": map[string]api.JSONSchema{
				"name": {
					"type": "string",
				},
				"age": {
					"type": "number",
				},
			},
		},
	}

	actual := JSONResponse(result, schema, log.NewFestive(0))
	assert.False(T, actual)
}

func TestJSONResponse_Array(T *testing.T) {
	result := &contract.OperationResult{
		ResponseBytes: []byte("[1, 2, 3, 4]"),
	}

	schema := &api.Schema{
		JSONSchema: api.JSONSchema{
			"type": "array",
			"items": api.JSONSchema{
				"type": "integer",
			},
		},
	}

	actual := JSONResponse(result, schema, log.NewFestive(0))
	assert.True(T, actual)
}

func TestJSONResponse_Array_False_Schema(T *testing.T) {
	result := &contract.OperationResult{
		ResponseBytes: []byte("[\"1\", 2, \"3\", 4]"),
	}

	schema := &api.Schema{
		JSONSchema: api.JSONSchema{
			"type": "array",
			"items": api.JSONSchema{
				"type": "integer",
			},
		},
	}

	actual := JSONResponse(result, schema, log.NewFestive(0))
	assert.False(T, actual)
}

func TestJSONResponse_Array_False_Unmarshal(T *testing.T) {
	result := &contract.OperationResult{
		ResponseBytes: []byte("[1, 2, \"3', 3RRR0RRR "),
	}

	schema := &api.Schema{
		JSONSchema: api.JSONSchema{
			"type": "array",
			"items": api.JSONSchema{
				"type": "integer",
			},
		},
	}

	actual := JSONResponse(result, schema, log.NewFestive(0))
	assert.False(T, actual)
}
