package test

import (
	"net/http"

	"github.com/x1n13y84issmd42/oasis/src/api"
	"github.com/x1n13y84issmd42/oasis/src/log"
)

// Operation performs a test of an operation by requesting a path
// and validating the received response headers & content against
// the definitions founds in an OAS spec file.
func Operation(specHost *api.Host, specOp *api.Operation, params *api.OperationParameters, logger log.ILogger) bool {
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

	req := specReq.CreateRequest(specHost)

	// Adding auth credentials.
	if specOp.Security != nil {
		specOp.Security.Secure(req)
	}

	// Requesting.
	logger.Requesting(req.Method, req.URL.String())
	response, err := client.Do(req)

	if err != nil {
		logger.Error(err)
		return false
	}

	// Actual testing starts here.
	res := Response(response, specResp, logger)

	if res {
		logger.OperationOK(specOp)
	} else {
		logger.OperationFail(specOp)
	}

	return res
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
