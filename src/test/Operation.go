package test

import (
	"net/http"

	"github.com/x1n13y84issmd42/oasis/src/api"
	"github.com/x1n13y84issmd42/oasis/src/log"
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
func Operation(specHost *api.Host, specOp *api.Operation, params *api.OperationParameters, logger log.ILogger) (result *OperationResult) {
	result = &OperationResult{}

	logger.TestingOperation(specOp)
	// Creating a request.
	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	specReq := SelectRequest(specOp, params)
	specResp := SelectResponse(specOp, params)
	//TODO: check if these != nil

	result.SpecRequest = specReq
	result.SpecResponse = specResp

	req := specReq.CreateRequest(specHost)
	result.HTTPRequest = req

	// Adding auth credentials.
	if specOp.Security != nil {
		specOp.Security.Secure(req)
	}

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
	result.Success = Response(response, specResp, logger)

	if result.Success {
		logger.OperationOK(specOp)
	} else {
		logger.OperationFail(specOp)
	}

	return
}

// SelectRequest selects a request data to use based on the params.Request.ContentTypeHint value.
func SelectRequest(specOp *api.Operation, params *api.OperationParameters) *api.Request {
	filterCT := func(specReq *api.Request) bool { return true }

	if params.Request.ContentTypeHint != "" {
		filterCT = func(specReq *api.Request) bool {
			return specReq.Headers.Get("Content-Type") == params.Request.ContentTypeHint
		}
	}

	for _, specReq := range specOp.Requests {
		if filterCT(specReq) {
			return specReq
		}
	}

	return nil
}

// SelectResponse selects response data to use based on the params.Response data.
func SelectResponse(specOp *api.Operation, params *api.OperationParameters) *api.Response {
	filterCT := func(specResp *api.Response) bool { return true }
	filterStatus := func(specResp *api.Response) bool { return true }

	if params.Response.ContentTypeHint != "" {
		filterCT = func(specResp *api.Response) bool {
			return specResp.ContentType == params.Response.ContentTypeHint
		}
	}

	if params.Response.StatusHint != 0 {
		filterCT = func(apiResp *api.Response) bool {
			return apiResp.StatusCode == params.Response.StatusHint
		}
	}

	for _, specResp := range specOp.Responses {
		if filterCT(specResp) && filterStatus(specResp) {
			return specResp
		}
	}

	return nil
}
