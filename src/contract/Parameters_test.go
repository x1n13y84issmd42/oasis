package contract_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/x1n13y84issmd42/oasis/src/contract"
	"github.com/x1n13y84issmd42/oasis/src/errors"
	"github.com/x1n13y84issmd42/oasis/src/params"
)

func Test_OperationData(T *testing.T) {
	ppp := contract.NewPOCTParameterProvider()

	srcP := params.NewMemoryParameterSource()
	srcO := params.NewMemoryParameterSource()
	srcC := params.NewMemoryParameterSource()
	srcT := params.NewMemoryParameterSource()

	ppp.AddSource(contract.ParameterSourceSpecPath, srcP)
	ppp.AddSource(contract.ParameterSourceSpecOp, srcO)
	ppp.AddSource(contract.ParameterSourceCLI, srcC)
	ppp.AddSource(contract.ParameterSourceTestOutput, srcT)

	srcP.Data["p_a"] = "p_a_path"
	srcP.Data["p_b"] = "p_b_path"
	srcP.Data["p_c"] = "p_c_path"
	srcP.Data["p_d"] = "p_d_path"

	srcO.Data["p_b"] = "p_b_op"
	srcO.Data["p_c"] = "p_c_op"
	srcO.Data["p_d"] = "p_d_op"

	srcC.Data["p_c"] = "p_c_cli"
	srcC.Data["p_d"] = "p_d_cli"

	srcT.Data["p_d"] = "p_d_test_out"

	T.Run("Require", func(T *testing.T) {

		ppp.Require("foo")
		ppp.Require("bar")
		ppp.Require("foo")
		ppp.Require("foo")
		ppp.Require("qeq")
		ppp.Require("foo")

		assert.Equal(T, 3, len(ppp.Required))
	})

	T.Run("Values", func(T *testing.T) {
		expected := map[string]string{
			"p_a": "p_a_path",
			"p_b": "p_b_op",
			"p_c": "p_c_cli",
			"p_d": "p_d_test_out",
		}

		assert.Equal(T, expected, ppp.Values())
	})

	T.Run("Iterate", func(T *testing.T) {
		expected := []string{
			"p_a_path",
			"p_b_op",
			"p_c_cli",
			"p_d_test_out",
		}

		actual := []string{}

		for pt := range ppp.Iterate() {
			actual = append(actual, pt[1])
		}

		assert.Equal(T, expected, actual)
	})

	T.Run("Get", func(T *testing.T) {
		expected := []string{
			"p_a_path",
			"p_b_op",
			"p_c_cli",
			"p_d_test_out",
		}

		actual := []string{
			ppp.Get("p_a"),
			ppp.Get("p_b"),
			ppp.Get("p_c"),
			ppp.Get("p_d"),
		}

		assert.Equal(T, expected, actual)
	})

	T.Run("Validate", func(T *testing.T) {
		ppp.Require("p_x")

		srcP.Data["qeq"] = "qeq_path"

		expected := errors.NoParameters([]string{"foo", "bar", "p_x"}, nil)

		assert.Equal(T, expected, ppp.Validate())
	})

	T.Run("Validate/NoError", func(T *testing.T) {
		srcP.Data["foo"] = "foo_path"
		srcP.Data["bar"] = "bar_path"
		srcP.Data["p_x"] = "p_x_path"

		assert.Nil(T, ppp.Validate())
	})
}
