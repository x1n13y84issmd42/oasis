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
		resolver := openapi3.NewDataResolver(log.NewPlain(0), spec.OAS, nil, nil)

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
		resolver := openapi3.NewDataResolver(log.NewPlain(0), spec.OAS, nil, nil)

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
		resolver := openapi3.NewDataResolver(log, spec.OAS, nil, nil)

		actual := resolver.Host("INVALID_HOST_NAME")

		expected := params.NoSource(errors.NotFound("Host", "IRRELEVANT", nil), log)

		assert.IsType(T, expected, actual)
	})

	T.Run("MetaData/OK", func(T *testing.T) {
		log := log.NewPlain(0)
		resolver := openapi3.NewDataResolver(log, spec.OAS, nil, &spec.OAS.Paths["/pet/{petId}"].Get.Responses)

		actualStatus, actualCT, actualSpecResp, actualSpecMT, err := resolver.MetaData(0, "")

		assert.Nil(T, err)
		assert.Equal(T, actualStatus, 200)
		assert.Equal(T, actualCT, "application/json")
		assert.NotNil(T, actualSpecResp)
		assert.NotNil(T, actualSpecMT)
	})

	T.Run("CollectHeaders", func(T *testing.T) {
		log := log.NewPlain(0)
		op := spec.OAS.Paths["/pet/{petId}"].Get
		resolver := openapi3.NewDataResolver(log, spec.OAS, nil, &op.Responses)

		actual, err := resolver.CollectHeaders(op.Responses["200"].Value)

		expected := []openapi3.ResolverExpectedHeader{
			{
				Name:     "X-Expires-After",
				Required: true,
			},

			{
				Name:     "X-Rate-Limit",
				Required: false,
			},
		}

		for i := range actual {
			actual[i].Schema = nil
		}

		assert.Nil(T, err)
		assert.Equal(T, expected, actual)
	})

	T.Run("MetaData/StatusError", func(T *testing.T) {
		log := log.NewPlain(0)
		resolver := openapi3.NewDataResolver(log, spec.OAS, nil, &spec.OAS.Paths["/pet/{petId}"].Get.Responses)

		expectedError := errors.NotFound("spec response", "201", nil)
		// Otherwise TheCaller points to this ^ place.
		expectedError.TheCaller = ""

		actualStatus, actualCT, actualSpecResp, actualSpecMT, err := resolver.MetaData(201, "")

		assert.IsType(T, expectedError, err)
		xerr := err.(errors.ErrNotFound)
		xerr.TheCaller = ""
		assert.Equal(T, expectedError, xerr)

		assert.Equal(T, actualStatus, 0)
		assert.Equal(T, actualCT, "")
		assert.Nil(T, actualSpecResp)
		assert.Nil(T, actualSpecMT)
	})

	T.Run("MetaData/CTError", func(T *testing.T) {
		log := log.NewPlain(0)
		resolver := openapi3.NewDataResolver(log, spec.OAS, nil, &spec.OAS.Paths["/pet/{petId}"].Get.Responses)

		expectedError := errors.NotFound("spec response", "image/png", nil)
		// Otherwise TheCaller points to this ^ place.
		expectedError.TheCaller = ""

		actualStatus, actualCT, actualSpecResp, actualSpecMT, err := resolver.MetaData(0, "image/png")

		assert.IsType(T, expectedError, err)
		xerr := err.(errors.ErrNotFound)
		xerr.TheCaller = ""
		assert.Equal(T, expectedError, xerr)

		assert.Equal(T, actualStatus, 0)
		assert.Equal(T, actualCT, "")
		assert.Nil(T, actualSpecResp)
		assert.Nil(T, actualSpecMT)
	})

	T.Run("MakeSchema/InvalidSchema/Marshal", func(T *testing.T) {
		resolver := openapi3.NewDataResolver(log.NewPlain(0), spec.OAS, nil, nil)

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

		resolver := openapi3.NewDataResolver(log.NewPlain(0), OAS, nil, nil)

		actual, actualErr := resolver.MakeSchema(schemaName, &schema)

		assert.Nil(T, actual)
		assert.Equal(T, expectedErr.Error(), actualErr.Error())
	})
}
