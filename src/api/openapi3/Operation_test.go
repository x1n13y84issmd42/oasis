package openapi3_test

import (
	"testing"

	kinopenapi3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/assert"
	"github.com/x1n13y84issmd42/oasis/src/api"
	"github.com/x1n13y84issmd42/oasis/src/api/openapi3"
	"github.com/x1n13y84issmd42/oasis/src/log"
)

func Test_Operation(T *testing.T) {
	log := log.NewPlain(0)
	oasPath := "/foo/{bar}"

	OAS := &kinopenapi3.Swagger{}

	oasOperation := &kinopenapi3.Operation{
		OperationID: "testop1",
		Summary:     "test op",
		Description: "a very test op",
		Responses: kinopenapi3.Responses{
			"200": &kinopenapi3.ResponseRef{
				Value: &kinopenapi3.Response{},
			},
		},
	}

	oasPathItem := &kinopenapi3.PathItem{
		Get: oasOperation,
	}

	OAS.Paths = kinopenapi3.Paths{}
	OAS.Paths[oasPath] = oasPathItem

	op := &openapi3.Operation{
		OperationPrototype: api.NewOperationPrototype(log),
		RequestMethod:      "OPTIONS",
		RequestPath:        oasPath,
		SpecOp:             oasOperation,
		SpecPath:           oasPathItem,
		Resolver:           openapi3.NewDataResolver(log, OAS, &oasOperation.Responses),
	}

	op.OperationPrototype.Operation = op

	T.Run("ID", func(T *testing.T) {
		assert.Equal(T, "testop1", op.ID())
	})

	T.Run("Name", func(T *testing.T) {
		assert.Equal(T, "test op", op.Name())
	})

	T.Run("Description", func(T *testing.T) {
		assert.Equal(T, "a very test op", op.Description())
	})

	T.Run("Method", func(T *testing.T) {
		assert.Equal(T, "OPTIONS", op.Method())
	})

	T.Run("Path", func(T *testing.T) {
		assert.Equal(T, oasPath, op.Path())
	})

	T.Run("Resolve", func(T *testing.T) {
		assert.NotNil(T, op.Resolve())
	})
}
