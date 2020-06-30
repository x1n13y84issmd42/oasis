package test

import (
	"github.com/x1n13y84issmd42/oasis/src/contract"
)

// Success creates a new OperationResult instance with
// the Success field set to true to assure maximum SUCCESS.
func Success() *contract.OperationResult {
	return &contract.OperationResult{
		Success: true,
	}
}

// Operation performs a test of an operation by requesting a path
// and validating the received response headers & content against
// the definitions founds in an OAS spec file.
func Operation(op contract.Operation, enrichment *[]contract.RequestEnrichment, v contract.Validator, log contract.Logger) *contract.OperationResult {
	// Creating a request.
	req := NewRequest(op, log)

	// Extending the request with stuff.
	for _, en := range *enrichment {
		req.Enrich(en)
	}

	// Requesting.
	result := req.Execute()

	// Testing & returning.
	result = v.Validate(result)

	if result.Success {
		log.OperationOK()
	} else {
		log.OperationFail()
	}

	return result
}
