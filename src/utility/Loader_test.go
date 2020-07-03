package utility_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/x1n13y84issmd42/oasis/src/api"
	"github.com/x1n13y84issmd42/oasis/src/api/openapi3"
	"github.com/x1n13y84issmd42/oasis/src/log"
	"github.com/x1n13y84issmd42/oasis/src/utility"
)

func Test_Loader(T *testing.T) {
	T.Run("OK", func(T *testing.T) {
		spec := utility.Load("../../spec/test/oas3_ok.yaml", log.NewPlain(1))
		_, ok := spec.(*openapi3.Spec)
		assert.True(T, ok)
	})

	T.Run("Failure", func(T *testing.T) {
		spec := utility.Load("A/VERY/WRONG/PATH.yaml", log.NewPlain(0))
		_, ok := spec.(api.NullSpec)
		assert.True(T, ok)
	})
}
