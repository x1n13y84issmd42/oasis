package openapi3_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/x1n13y84issmd42/oasis/src/log"
	"github.com/x1n13y84issmd42/oasis/src/params"
	"github.com/x1n13y84issmd42/oasis/src/utility"
)

func Test_Spec(T *testing.T) {
	log := log.NewPlain(0)
	spec := utility.Load("../../../spec/test/oas3.yaml", log)

	T.Run("Info", func(T *testing.T) {
		assert.Equal(T, "Swagger Petstore Test Version", spec.Title())
		assert.Equal(T, "This is a test version of the Swagger Petstore API.", spec.Description())
		assert.Equal(T, "1.0.0", spec.Version())
	})

	T.Run("Operations", func(T *testing.T) {
		expected := []string{
			"getPetById",
			"getUserByName",
			"updateUser",
			"deleteUser",
		}

		actual := []string{}

		for op := range spec.Operations() {
			actual = append(actual, op.ID())
		}

		assert.Equal(T, expected, actual)
	})

	T.Run("GetOperation", func(T *testing.T) {
		op := spec.GetOperation("deleteUser")

		pp := params.NewMemorySource("")
		pp.Add("username", "GOPHER")
		op.Data().URL.Load(pp)
		op.Data().URL.Load(op.Resolve().Host(""))

		assert.Equal(T, "https://petstore.swagger.io/v2/user/GOPHER", op.Data().URL.String())
	})
}
