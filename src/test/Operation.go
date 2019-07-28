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

// Run performs a test.
func (op Operation) Run(requestContentType string, responseStatus int, responseContentType string) {
	op.Log.TestingOperation(op.Operation)

	client := http.Client{}
	req := op.makeRequest(requestContentType)

	response, err := client.Do(req)

	if err != nil {
		fmt.Println("Something happened: ", err)
		return
	}

	resp := op.getResponse(responseStatus, responseContentType)
	if resp == nil {
		fmt.Printf("No response for Status of %d & Content-Type of \"%s\"\n", responseStatus, responseContentType)
		return
	}

	tResp := NewResponse(resp, op.Log)
	tResp.Test(response)
}

func (op Operation) makeRequest(CT string) *http.Request {
	URL := fmt.Sprintf("%s%s", op.Host.URL, op.Operation.Path)
	fmt.Printf("Requesting %s\n", URL)
	req, _ := http.NewRequest(op.Operation.Method, URL, nil)
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
