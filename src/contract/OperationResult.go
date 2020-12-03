package contract

import (
	"net/http"
)

// OperationResult describes the outcome of an operation test.
// It is bo te used as a possible source of data for subsequent tests.
type OperationResult struct {
	Success       bool
	HTTPRequest   *http.Request
	HTTPResponse  *http.Response
	ResponseBytes []byte
}

// And creates a new OperationResult instance with the Success field assigned
// the result of r1.Success && r2.Success.
func (r1 *OperationResult) And(r2 *OperationResult) *OperationResult {
	return &OperationResult{
		Success: r1.Success && r2.Success,
	}
}

// OperationResults is a map of operation results.
type OperationResults map[string]*OperationResult
