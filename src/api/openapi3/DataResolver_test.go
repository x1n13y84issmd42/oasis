package openapi3_test

import (
	"testing"

	kinopenapi3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/assert"
	"github.com/x1n13y84issmd42/oasis/src/api/openapi3"
	"github.com/x1n13y84issmd42/oasis/src/errors"
	"github.com/x1n13y84issmd42/oasis/src/log"
	"github.com/x1n13y84issmd42/oasis/src/params"
)

func Test_DataResolver(T *testing.T) {
	spec, _ := openapi3.Load("../../../spec/test/oas3.yaml", log.NewPlain(1))

	T.Run("Host", func(T *testing.T) {
		resolver := openapi3.NewDataResolver(log.NewPlain(0), spec.OAS, nil)

		src := resolver.Host("")

		expected := []string{
			"https://petstore.swagger.io/v2",
		}

		actual := []string{}

		for p := range src.Iterate() {
			actual = append(actual, p.V())
		}

		assert.Equal(T, expected, actual)
	})

	T.Run("Host/Named", func(T *testing.T) {
		resolver := openapi3.NewDataResolver(log.NewPlain(0), spec.OAS, nil)

		src := resolver.Host("HTTP")

		expected := []string{
			"http://petstore.swagger.io/v2",
		}

		actual := []string{}

		for p := range src.Iterate() {
			actual = append(actual, p.V())
		}

		assert.Equal(T, expected, actual)
	})

	T.Run("Host/NoSource", func(T *testing.T) {
		log := log.NewPlain(0)
		resolver := openapi3.NewDataResolver(log, spec.OAS, nil)

		actual := resolver.Host("INVALID_HOST_NAME")

		expected := params.NoSource(errors.NotFound("Host", "IRRELEVANT", nil), log)

		assert.IsType(T, expected, actual)
	})

	T.Run("MakeSchema/InvalidSchema/Marshal", func(T *testing.T) {
		resolver := openapi3.NewDataResolver(log.NewPlain(0), spec.OAS, nil)

		schemaName := "schema_one"
		expectedErr := errors.InvalidSchema(schemaName, "Failed to marshal the schema.", nil)

		schema := kinopenapi3.Schema{
			Type:    "object",
			Example: make(chan int),
		}

		actual, actualErr := resolver.MakeSchema(schemaName, &schema)

		assert.Nil(T, actual)
		assert.Equal(T, expectedErr.Error(), actualErr.Error())
	})

	T.Run("MakeSchema/InvalidSchema/MarshalComponents", func(T *testing.T) {
		schemaName := "schema_one"
		expectedErr := errors.InvalidSchema(schemaName, "Failed to marshal Components.", nil)

		schema := kinopenapi3.Schema{
			Type: "object",
		}

		OAS := &kinopenapi3.Swagger{
			Components: kinopenapi3.Components{
				Schemas: map[string]*kinopenapi3.SchemaRef{
					"Errorneous component": {
						Value: &kinopenapi3.Schema{
							Example: make(chan int),
						},
					},
				},
			},
		}

		resolver := openapi3.NewDataResolver(log.NewPlain(0), OAS, nil)

		actual, actualErr := resolver.MakeSchema(schemaName, &schema)

		assert.Nil(T, actual)
		assert.Equal(T, expectedErr.Error(), actualErr.Error())
	})
}
