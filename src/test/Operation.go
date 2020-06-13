package test

import (
	"net/http"

	"github.com/x1n13y84issmd42/oasis/src/api"
	"github.com/x1n13y84issmd42/oasis/src/contract"
)

// OperationResult describes the outcome of an operation test.
// It is bo te used as a possible source of data for subsequent tests.
type OperationResult struct {
	Success      bool
	SpecRequest  *api.Request
	HTTPRequest  *http.Request
	SpecResponse *api.Response
	HTTPResponse *http.Response
}

// And creates a new OperationResult instance with the Success field assigned
// the result of r1.Success && r2.Success.
func (r1 *OperationResult) And(r2 *OperationResult) *OperationResult {
	return &OperationResult{
		Success: r1.Success && r2.Success,
	}
}

// Success creates a new OperationResult instance with
// the Success field set to true to assure maximum SUCCESS.
func Success() *OperationResult {
	return &OperationResult{
		Success: true,
	}
}

// Operation performs a test of an operation by requesting a path
// and validating the received response headers & content against
// the definitions founds in an OAS spec file.
func Operation(specHost string, specOp contract.Operation, params *contract.OperationParameters, logger contract.Logger) (result *OperationResult) {
	result = &OperationResult{}

	logger.TestingOperation(specOp)
	// Creating a request.
	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	req := specOp.CreateRequest()
	result.HTTPRequest = req

	// // Adding auth credentials.
	// if specOp.Security != nil {
	// 	specOp.Security.Secure(req)
	// }

	// Requesting.
	logger.Requesting(req.Method, req.URL.String())
	response, err := client.Do(req)
	result.HTTPResponse = response

	if err != nil {
		logger.Error(err)
		result.Success = false
		return
	}

	// Actual testing starts here.
	// result.Success = Response(response, specResp, logger)

	if result.Success {
		logger.OperationOK(specOp)
	} else {
		logger.OperationFail(specOp)
	}

	return
}
