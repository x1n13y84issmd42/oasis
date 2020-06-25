package params_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/x1n13y84issmd42/oasis/src/errors"
	"github.com/x1n13y84issmd42/oasis/src/params"
)

func Test_MultiSet(T *testing.T) {
	T.Run("Load&Iterate", func(T *testing.T) {
		src1 := params.NewMemorySource()
		src1.Add("A", "The aye")
		src1.Add("B", "The bee")
		src1.Add("C", "The sea")
		src1.Add("D", "The D")

		src2 := params.NewMemorySource()
		src2.Add("A", "The aye aye")
		src2.Add("D", "The D #2")

		set := params.NewMultiSet()
		set.Load(src1)
		set.Load(src2)

		expected := []string{
			"The aye",
			"The aye aye",
			"The bee",
			"The sea",
			"The D",
			"The D #2",
		}

		actual := []string{}

		for p := range set.Iterate() {
			actual = append(actual, p.V)
		}

		assert.Equal(T, expected, actual)
	})

	T.Run("Require&Validate Fail", func(T *testing.T) {
		src := params.NewMemorySource()
		src.Add("A", "The aye")
		src.Add("B", "The bee")

		set := params.NewMultiSet()
		set.Load(src)

		set.Require("A")
		set.Require("C")
		set.Require("D")
		set.Require("D")
		set.Require("D")

		expected := errors.NoParameters([]string{"C", "D"}, nil)

		assert.Equal(T, expected, set.Validate())
	})

	T.Run("Require&Validate Success", func(T *testing.T) {
		src := params.NewMemorySource()
		src.Add("A", "The aye")
		src.Add("B", "The bee")

		set := params.NewMultiSet()
		set.Load(src)

		set.Require("A")
		set.Require("B")
		set.Require("B")
		set.Require("B")
		set.Require("B")

		assert.Nil(T, set.Validate())
	})
}
