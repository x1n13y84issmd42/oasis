package api_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/x1n13y84issmd42/oasis/src/api"
	"github.com/x1n13y84issmd42/oasis/src/log"
	"github.com/x1n13y84issmd42/oasis/src/params"
)

// testOp provides access to OAS3-specific API data.
type testOp struct {
	*api.OperationPrototype
}

// ID ...
func (op *testOp) ID() string {
	return "testop"
}

// Name ...
func (op *testOp) Name() string {
	return "Test Op"
}

// Description ...
func (op *testOp) Description() string {
	return "A Test Operation"
}

// Method ...
func (op *testOp) Method() string {
	return "OPTIONS"
}

// Path ...
func (op *testOp) Path() string {
	return "/foo/{bar}"
}

func Test_OperationPrototype(T *testing.T) {
	log := log.NewPlain(0)

	op := &testOp{
		OperationPrototype: api.NewOperationPrototype(log),
	}

	op.OperationPrototype.Operation = op

	src := params.NewMemorySource("test")
	src.Add(params.KeyHost, "example.com")
	src.Add("bar", "B4R_VALUE")

	op.OperationPrototype.Data().URL = params.URL(op.Path(), log)
	op.OperationPrototype.Data().URL.Load(src)

	T.Run("GetRequest", func(T *testing.T) {
		expectedURL := "example.com/foo/B4R_VALUE"
		expectedMethod := "OPTIONS"
		req, _ := op.OperationPrototype.GetRequest()

		assert.Equal(T, expectedURL, req.URL.String())
		assert.Equal(T, expectedMethod, req.Method)
	})
}
