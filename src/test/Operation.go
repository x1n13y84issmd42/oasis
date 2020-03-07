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
func (test Operation) Run(requestContentType string, responseStatus int, responseContentType string) bool {
	test.Log.TestingOperation(test.Operation)

	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req := test.makeRequest(requestContentType)

	response, err := client.Do(req)

	if err != nil {
		test.Log.Error(err)
		return false
	}

	apiResp := test.getResponse(responseStatus, responseContentType)
	if apiResp == nil {
		test.Log.ResponseNotFound(responseContentType, responseStatus)
		return false
	}

	//	Testing the response.
	tResp := NewResponse(apiResp, test.Log)
	return tResp.Test(response)
}

func (test Operation) makeRequest(CT string) *http.Request {
	URL := fmt.Sprintf("%s%s", test.Host.URL, test.Operation.Path)
	test.Log.Requesting(URL)
	req, _ := http.NewRequest(test.Operation.Method, URL, nil)
	//TODO: req body & CT
	return req
}

func (test Operation) getResponse(status int, CT string) *api.Response {
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

	for _, resp := range *test.Operation.Responses {
		if filterStatus(resp) && filterCT(resp) {
			return &resp
		}
	}

	return nil
}
