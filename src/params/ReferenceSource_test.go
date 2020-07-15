package params_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/x1n13y84issmd42/oasis/src/contract"
	"github.com/x1n13y84issmd42/oasis/src/log"
	"github.com/x1n13y84issmd42/oasis/src/params"
)

func Test_ReferenceSource(T *testing.T) {
	src := params.NewReferenceSource(log.NewPlain(0))

	result1 := &contract.OperationResult{
		ResponseBytes: []byte(`[{"name":"john"}]`),
	}

	result2 := &contract.OperationResult{
		ResponseBytes: []byte(`[{"values":[41, 42, 43]}]`),
	}

	src.AddReference("foo", "testMe", result1, "[0].name")
	src.AddReference("snek", "noStepOnSnek", result2, "[0].values[1]")

	expected := []string{
		"testMe.foo = john",
		"noStepOnSnek.snek = 42",
	}

	actual := []string{}

	for p := range src.Iterate() {
		actual = append(actual, fmt.Sprintf("%s.%s = %s", p.Source, p.N, p.V()))
	}

	assert.Equal(T, expected, actual)
}
