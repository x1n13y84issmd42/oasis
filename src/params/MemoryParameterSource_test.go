package params_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/x1n13y84issmd42/oasis/src/params"
)

func Test_MemoryParameterSource(T *testing.T) {
	src := params.NewMemoryParameterSource()
	src.Data["b"] = "BB"
	src.Data["a"] = "AAAAA"

	T.Run("Get", func(T *testing.T) {
		assert.Equal(T, "AAAAA", src.Get("a"))
		assert.Equal(T, "BB", src.Get("b"))
		assert.Equal(T, "", src.Get("c"))
	})

	T.Run("Iterate", func(T *testing.T) {
		expected := []string{
			"AAAAA",
			"BB",
		}

		actual := []string{}

		for pt := range src.Iterate() {
			actual = append(actual, pt.V)
		}

		assert.Equal(T, expected, actual)
	})
}
