package openapi3

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/assert"
	"github.com/x1n13y84issmd42/oasis/src/api"
	"github.com/x1n13y84issmd42/oasis/src/log"
)

func TestCreatePath(T *testing.T) {
	path := "/foo/{param1}/bar/{param2}/{param3}"
	expectedPath := "/foo/p1_from_op/bar/p2_from_path/p3_from_override"

	OAS := &openapi3.Swagger{}
	oasOperation := openapi3.Operation{
		Parameters: openapi3.Parameters{
			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					In:       "path",
					Name:     "param1",
					Required: true,
					Example:  "p1_from_op",
				},
			},

			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					In:      "query",
					Name:    "param2",
					Example: "p2_from_op",
				},
			},

			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					In:       "path",
					Name:     "param3",
					Required: true,
					Example:  "p3_from_op",
				},
			},
		},
	}
	oasPathItem := &openapi3.PathItem{
		Get: &oasOperation,
		Parameters: openapi3.Parameters{
			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					In:       "path",
					Name:     "param1",
					Required: true,
					Example:  "p1_from_path",
				},
			},

			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					In:       "path",
					Name:     "param2",
					Required: true,
					Example:  "p2_from_path",
				},
			},
		},
	}
	OAS.Paths = openapi3.Paths{}
	OAS.Paths[path] = oasPathItem

	params := api.OperationParameters{
		Path: api.PathParameters{
			"param1": "",
			"param3": "p3_from_override",
		},
	}

	spec := Spec{
		Log: log.NewFestive(0),
		OAS: OAS,
	}

	actualPath, err := spec.CreatePath(path, oasPathItem, &oasOperation, &params)

	assert.Equal(T, actualPath, expectedPath)
	assert.Nil(T, err)
}
