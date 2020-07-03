package api_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/x1n13y84issmd42/oasis/src/api"
)

func Test_Schema(T *testing.T) {
	T.Run("Cast integer", func(T *testing.T) {
		schema := &api.Schema{
			JSONSchema: api.JSONSchema{
				"type": "integer",
			},
		}

		assert.Equal(T, int64(42), schema.Cast("42"))
	})

	T.Run("Cast number", func(T *testing.T) {
		schema := &api.Schema{
			JSONSchema: api.JSONSchema{
				"type": "number",
			},
		}

		assert.Equal(T, 42.25, schema.Cast("42.25"))
	})

	T.Run("Cast boolean", func(T *testing.T) {
		schema := &api.Schema{
			JSONSchema: api.JSONSchema{
				"type": "boolean",
			},
		}

		assert.Equal(T, true, schema.Cast("true"))
	})

	T.Run("Cast passthrough", func(T *testing.T) {
		schema := &api.Schema{
			JSONSchema: api.JSONSchema{
				"type": "boolean",
			},
		}

		assert.Equal(T, "foobar", schema.Cast("foobar"))
	})
}
