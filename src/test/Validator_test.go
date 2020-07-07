package test_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/x1n13y84issmd42/oasis/src/contract"
	"github.com/x1n13y84issmd42/oasis/src/log"
	"github.com/x1n13y84issmd42/oasis/src/test"
)

func Test_Validator(T *testing.T) {
	log := log.NewPlain(0)
	v := test.NewValidator(log)

	v.Expect(func(req *http.Response) bool {
		return false
	})

	r := &contract.OperationResult{
		Success: true,
	}

	r = v.Validate(r)

	assert.False(T, r.Success)
}
