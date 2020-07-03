package openapi3_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/x1n13y84issmd42/oasis/src/api/openapi3"
	"github.com/x1n13y84issmd42/oasis/src/log"
)

func Test_Loader(T *testing.T) {
	T.Run("OK", func(T *testing.T) {
		spec, specerr := openapi3.Load("../../../spec/test/oas3_ok.yaml", log.NewPlain(1))
		assert.NotNil(T, spec)
		assert.Nil(T, specerr)
	})

	T.Run("Failure", func(T *testing.T) {
		spec, specerr := openapi3.Load("A/VERY/WRONG/PATH.yaml", log.NewPlain(0))
		assert.Nil(T, spec)
		assert.NotNil(T, specerr)
	})
}
