package params_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/x1n13y84issmd42/oasis/src/params"
)

func Test_Set(T *testing.T) {
	T.Run("Load&Iterate", func(T *testing.T) {
		src1 := params.NewMemorySource("test_1")
		src1.Add("A", "The aye")
		src1.Add("B", "The bee")
		src1.Add("C", "The sea")
		src1.Add("D", "The D")

		src2 := params.NewMemorySource("7357**")
		src2.Add("A", "The aye aye")
		src2.Add("D", "The D #2")

		set := params.NewSet("yolo")
		set.Load(src1)
		set.Load(src2)

		expected := []string{
			"A 7357** The aye aye",
			"B test_1 The bee",
			"C test_1 The sea",
			"D 7357** The D #2",
		}

		actual := []string{}

		for p := range set.Iterate() {
			actual = append(actual, p.N+" "+p.Source+" "+p.V())
		}

		assert.Equal(T, expected, actual)
	})
}
