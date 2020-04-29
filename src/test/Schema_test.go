package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/x1n13y84issmd42/oasis/src/api"
	"github.com/x1n13y84issmd42/oasis/src/log"
)

func TestSchema_False_Invalid(T *testing.T) {
	specSchema := &api.Schema{
		JSONSchema: api.JSONSchema{
			"type": "integer",
		},
	}

	data := "string"
	actual := Schema(data, specSchema, log.NewFestive(0))
	assert.False(T, actual)
}

func TestSchema_False_Error(T *testing.T) {
	specSchema := &api.Schema{
		JSONSchema: api.JSONSchema{
			"type": "integer",
		},
	}

	data := make(chan int)
	actual := Schema(data, specSchema, log.NewFestive(0))
	assert.False(T, actual)
}

func TestSchema(T *testing.T) {
	specSchema := &api.Schema{
		JSONSchema: api.JSONSchema{
			"type": "integer",
		},
	}

	data := 42
	actual := Schema(data, specSchema, log.NewFestive(0))
	assert.True(T, actual)
}
