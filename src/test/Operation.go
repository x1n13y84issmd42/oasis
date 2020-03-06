package test

import (
	"fmt"
	"net/http"

	"github.com/x1n13y84issmd42/goasis/src/api"
	"github.com/x1n13y84issmd42/goasis/src/log"
)

// Operation performs a test of an operation by requesting a path
// and validating thereceived response headers & content.
type Operation struct {
	Log       log.ILogger
	Host      *api.Host
	Operation *api.Operation
}

// Run performs a test of an operation.
func (op Operation) Run(requestContentType string, responseStatus int, responseContentType string) bool {
	op.Log.TestingOperation(op.Operation)

	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req := op.makeRequest(requestContentType)

	response, err := client.Do(req)

	if err != nil {
		op.Log.Error(err)
		return false
	}

	resp := op.getResponse(responseStatus, responseContentType)
	if resp == nil {
		op.Log.ResponseNotFound(responseContentType, responseStatus)
		return false
	}

	//	Testing the response.
	tResp := NewResponse(resp, op.Log)
	return tResp.Test(response)
}

func (op Operation) makeRequest(CT string) *http.Request {
	URL := fmt.Sprintf("%s%s", op.Host.URL, op.Operation.Path)
	op.Log.Requesting(URL)
	req, _ := http.NewRequest(op.Operation.Method, URL, nil)
	//TODO: req body & CT
	return req
}

func (op Operation) getResponse(status int, CT string) *api.Response {
	filterCT := func(apiResp api.Response) bool {
		return apiResp.ContentType == CT
	}

	filterStatus := func(apiResp api.Response) bool {
		return apiResp.StatusCode == status
	}

	if CT == "*" {
		filterCT = func(api.Response) bool { return true }
	}

	if status == 0 {
		filterStatus = func(api.Response) bool { return true }
	}

	for _, resp := range *op.Operation.Responses {
		if filterStatus(resp) && filterCT(resp) {
			return &resp
		}
	}

	return nil
}
