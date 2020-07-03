package contract_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/x1n13y84issmd42/oasis/src/contract"
)

func Test_OperationResult(T *testing.T) {
	resT := &contract.OperationResult{
		Success: true,
	}

	resF := &contract.OperationResult{
		Success: false,
	}

	assert.True(T, resT.And(resT).Success)
	assert.False(T, resF.And(resT).Success)
	assert.False(T, resT.And(resF).Success)
}
