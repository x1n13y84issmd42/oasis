package test

import (
	"fmt"
	"net/http"

	"github.com/x1n13y84issmd42/goasis/src/api"
	"github.com/x1n13y84issmd42/goasis/src/log"
	"github.com/x1n13y84issmd42/goasis/src/test/security"
)

// Operation performs a test of an operation by requesting a path
// and validating thereceived response headers & content.
type Operation struct {
	Log       log.ILogger
	Host      *api.Host
	Operation *api.Operation
}

// Run performs a test of an operation by making a requests to the operation URL with the operation method
// and chosen Content-Type.
func (test Operation) Run(requestContentType string, responseStatus int, responseContentType string) bool {
	test.Log.TestingOperation(test.Operation)

	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req := test.createRequest(requestContentType)

	if test.Operation.Security != nil {
		test.Log.UsingSecurity(test.Operation.Security)
		security.NewSecurity(test.Operation.Security, test.Log).Secure(req)
	}

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

func (test Operation) createRequest(CT string) *http.Request {
	URL := fmt.Sprintf("%s%s", test.Host.URL, test.Operation.Path)
	test.Log.Requesting(URL)
	req, _ := http.NewRequest(test.Operation.Method, URL, nil)
	//TODO: use req body & CT when applicable
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
