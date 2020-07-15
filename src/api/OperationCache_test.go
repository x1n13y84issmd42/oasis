package api_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/x1n13y84issmd42/oasis/src/api"
	"github.com/x1n13y84issmd42/oasis/src/contract"
)

type op struct{}

// GetRequest reports an error.
func (op op) GetRequest() (*http.Request, error) {
	return nil, nil
}

// ID reports an error.
func (op op) ID() string {
	return ""
}

// Name reports an error.
func (op op) Name() string {
	return ""
}

// Description reports an error.
func (op op) Description() string {
	return ""
}

// Method reports an error.
func (op op) Method() string {
	return ""
}

// Path reports an error.
func (op op) Path() string {
	return ""
}

// Data reports an error.
func (op op) Data() *contract.OperationData {
	return nil
}

// Resolve reports an error.
func (op op) Resolve() contract.DataResolver {
	return nil
}

// Result reports an error.
func (op op) Result() *contract.OperationResult {
	return nil
}

type opaccess struct {
	RequestedOpID string
}

func (oa *opaccess) GetOperation(id string) contract.Operation {
	oa.RequestedOpID = id
	return nil
}

func (oa *opaccess) Operations() contract.OperationIterator {
	return nil
}

func Test_OperationAccess(T *testing.T) {
	access := &opaccess{}
	cache := api.NewOperationCache(access)

	opID := "theGloriousOp"

	cache.GetOperation(opID)
	assert.Equal(T, opID, access.RequestedOpID)

	access.RequestedOpID = ""
	cache.GetOperation(opID)
	assert.Equal(T, "", access.RequestedOpID)
}
