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
func Operation(op contract.Operation, enrichment *[]contract.RequestEnrichment, logger contract.Logger) (result *contract.OperationResult) {
	logger.TestingOperation(op)

	// Creating a request.
	req := NewRequest(op, logger)

	// Extending the request with stuff.
	for _, en := range *enrichment {
		req.Enrich(en)
	}

	// Requesting & returning.
	return req.Execute()
}
