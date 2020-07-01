package test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/x1n13y84issmd42/oasis/src/api"
	"github.com/x1n13y84issmd42/oasis/src/log"
)

func TestJSONResponse_String(T *testing.T) {
	httpResp := &http.Response{
		Body: ioutil.NopCloser(bytes.NewReader([]byte("\"anything works until it's quoted\""))),
	}

	specResp := &api.Response{
		Schema: &api.Schema{
			JSONSchema: api.JSONSchema{
				"type": "string",
			},
		},
	}

	actual := JSONResponse(httpResp, specResp.Schema, log.NewFestive(0))
	assert.True(T, actual)
}

func TestJSONResponse_String_False(T *testing.T) {
	httpResp := &http.Response{
		Body: ioutil.NopCloser(bytes.NewReader([]byte("unquoted malformed string\""))),
	}

	specResp := &api.Response{
		Schema: &api.Schema{
			JSONSchema: api.JSONSchema{
				"type": "string",
			},
		},
	}

	actual := JSONResponse(httpResp, specResp.Schema, log.NewFestive(0))
	assert.False(T, actual)
}

func TestJSONResponse_Number(T *testing.T) {
	httpResp := &http.Response{
		Body: ioutil.NopCloser(bytes.NewReader([]byte("42"))),
	}

	specResp := &api.Response{
		Schema: &api.Schema{
			JSONSchema: api.JSONSchema{
				"type": "number",
			},
		},
	}

	actual := JSONResponse(httpResp, specResp.Schema, log.NewFestive(0))
	assert.True(T, actual)
}

func TestJSONResponse_Number_False(T *testing.T) {
	httpResp := &http.Response{
		Body: ioutil.NopCloser(bytes.NewReader([]byte("not a number, not even close"))),
	}

	specResp := &api.Response{
		Schema: &api.Schema{
			JSONSchema: api.JSONSchema{
				"type": "number",
			},
		},
	}

	actual := JSONResponse(httpResp, specResp.Schema, log.NewFestive(0))
	assert.False(T, actual)
}

func TestJSONResponse_Boolean(T *testing.T) {
	httpResp := &http.Response{
		Body: ioutil.NopCloser(bytes.NewReader([]byte("false"))),
	}

	specResp := &api.Response{
		Schema: &api.Schema{
			JSONSchema: api.JSONSchema{
				"type": "boolean",
			},
		},
	}

	actual := JSONResponse(httpResp, specResp.Schema, log.NewFestive(0))
	assert.True(T, actual)
}

func TestJSONResponse_Boolean_False(T *testing.T) {
	httpResp := &http.Response{
		Body: ioutil.NopCloser(bytes.NewReader([]byte("truth"))),
	}

	specResp := &api.Response{
		Schema: &api.Schema{
			JSONSchema: api.JSONSchema{
				"type": "boolean",
			},
		},
	}

	actual := JSONResponse(httpResp, specResp.Schema, log.NewFestive(0))
	assert.False(T, actual)
}

func TestJSONResponse_Object(T *testing.T) {
	httpResp := &http.Response{
		Body: ioutil.NopCloser(bytes.NewReader([]byte("{\"name\":\"johnny\",\"age\":42}"))),
	}

	specResp := &api.Response{
		Schema: &api.Schema{
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
		},
	}

	actual := JSONResponse(httpResp, specResp.Schema, log.NewFestive(0))
	assert.True(T, actual)
}

func TestJSONResponse_Object_False_Schema(T *testing.T) {
	httpResp := &http.Response{
		Body: ioutil.NopCloser(bytes.NewReader([]byte("{\"name\":false,\"age\":\"42\"}"))),
	}

	specResp := &api.Response{
		Schema: &api.Schema{
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
		},
	}

	actual := JSONResponse(httpResp, specResp.Schema, log.NewFestive(0))
	assert.False(T, actual)
}

func TestJSONResponse_Object_False_Unmarshal(T *testing.T) {
	httpResp := &http.Response{
		Body: ioutil.NopCloser(bytes.NewReader([]byte("... INVALID JSON }"))),
	}

	specResp := &api.Response{
		Schema: &api.Schema{
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
		},
	}

	actual := JSONResponse(httpResp, specResp.Schema, log.NewFestive(0))
	assert.False(T, actual)
}

func TestJSONResponse_Array(T *testing.T) {
	httpResp := &http.Response{
		Body: ioutil.NopCloser(bytes.NewReader([]byte("[1, 2, 3, 4]"))),
	}

	specResp := &api.Response{
		Schema: &api.Schema{
			JSONSchema: api.JSONSchema{
				"type": "array",
				"items": api.JSONSchema{
					"type": "integer",
				},
			},
		},
	}

	actual := JSONResponse(httpResp, specResp.Schema, log.NewFestive(0))
	assert.True(T, actual)
}

func TestJSONResponse_Array_False_Schema(T *testing.T) {
	httpResp := &http.Response{
		Body: ioutil.NopCloser(bytes.NewReader([]byte("[\"1\", 2, \"3\", 4]"))),
	}

	specResp := &api.Response{
		Schema: &api.Schema{
			JSONSchema: api.JSONSchema{
				"type": "array",
				"items": api.JSONSchema{
					"type": "integer",
				},
			},
		},
	}

	actual := JSONResponse(httpResp, specResp.Schema, log.NewFestive(0))
	assert.False(T, actual)
}

func TestJSONResponse_Array_False_Unmarshal(T *testing.T) {
	httpResp := &http.Response{
		Body: ioutil.NopCloser(bytes.NewReader([]byte("[1, 2, \"3', 3RRR0RRR "))),
	}

	specResp := &api.Response{
		Schema: &api.Schema{
			JSONSchema: api.JSONSchema{
				"type": "array",
				"items": api.JSONSchema{
					"type": "integer",
				},
			},
		},
	}

	actual := JSONResponse(httpResp, specResp.Schema, log.NewFestive(0))
	assert.False(T, actual)
}
