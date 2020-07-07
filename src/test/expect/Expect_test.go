package expect_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/x1n13y84issmd42/oasis/src/api"
	"github.com/x1n13y84issmd42/oasis/src/log"
	"github.com/x1n13y84issmd42/oasis/src/test/expect"
)

func Test_Status(T *testing.T) {
	log := log.New("plain", 0)
	resp := &http.Response{
		StatusCode: 400,
	}

	T.Run("True", func(T *testing.T) {
		assert.True(T, expect.Status(400, log)(resp))
	})

	T.Run("False", func(T *testing.T) {
		assert.False(T, expect.Status(200, log)(resp))
	})
}

func Test_HeaderRequired(T *testing.T) {
	log := log.New("plain", 0)
	resp := &http.Response{
		Header: http.Header{
			"Content-Type": []string{
				"application/json",
			},
		},
	}

	T.Run("True", func(T *testing.T) {
		assert.True(T, expect.HeaderRequired("CoNtEnt-TyPe", log)(resp))
	})

	T.Run("False", func(T *testing.T) {
		assert.False(T, expect.HeaderRequired("x-thing", log)(resp))
	})
}

func Test_ContentType(T *testing.T) {
	log := log.New("plain", 0)
	resp := &http.Response{
		Header: http.Header{
			"Content-Type": []string{
				"application/json; charset=utf-8",
			},
		},
	}

	T.Run("True", func(T *testing.T) {
		assert.True(T, expect.ContentType("application/json", log)(resp))
	})

	T.Run("False", func(T *testing.T) {
		assert.False(T, expect.ContentType("text/html", log)(resp))
	})
}

func Test_HeaderSchema(T *testing.T) {
	log := log.New("plain", 0)
	resp := &http.Response{
		Header: http.Header{
			"X-Thing": []string{
				"456789",
			},
		},
	}

	T.Run("True", func(T *testing.T) {
		schema := &api.Schema{
			JSONSchema: api.JSONSchema{
				"type": "integer",
			},
		}

		assert.True(T, expect.HeaderSchema("X-Thing", schema, log)(resp))
	})

	T.Run("False", func(T *testing.T) {
		schema := &api.Schema{
			JSONSchema: api.JSONSchema{
				"type": "boolean",
			},
		}

		assert.False(T, expect.HeaderSchema("X-Thing", schema, log)(resp))
	})
}

func Test_ContentSchema(T *testing.T) {
	log := log.New("plain", 0)
	resp := &http.Response{
		StatusCode: 200,
		Header: http.Header{
			"Content-Type": []string{
				"application/json",
			},
			"X-Index": []string{
				"3",
			},
		},
		Body: ioutil.NopCloser(bytes.NewReader([]byte("42"))),
	}

	T.Run("True", func(T *testing.T) {
		schema := &api.Schema{
			JSONSchema: api.JSONSchema{
				"type": "integer",
			},
		}

		assert.True(T, expect.ContentSchema(schema, log)(resp))
	})

	T.Run("False", func(T *testing.T) {
		schema := &api.Schema{
			JSONSchema: api.JSONSchema{
				"type": "boolean",
			},
		}

		assert.False(T, expect.ContentSchema(schema, log)(resp))
	})
}
