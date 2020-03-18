package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/x1n13y84issmd42/goasis/src/api"
	"github.com/x1n13y84issmd42/goasis/src/log"
	"github.com/x1n13y84issmd42/goasis/src/test/security"
)

// Operation performs a test of an operation by requesting a path
// and validating the received response headers & content against
// the definitions founds in an OAS spec file.
type Operation struct {
	Log       log.ILogger
	Host      *api.Host
	Operation *api.Operation
}

// Run performs a test of an operation by making a requests to the operation URL with the operation method
// and chosen Content-Type.
func (test Operation) Run(requestContentType string, responseStatus int, responseContentType string) bool {
	test.Log.TestingOperation(test.Operation)

	// Creating a request.
	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req := test.createRequest(requestContentType)

	// Applying a security.
	if test.Operation.Security != nil {
		test.Log.UsingSecurity(test.Operation.Security)
		security.NewSecurity(test.Operation.Security, test.Log).Secure(req)
	}

	// Requesting.
	test.Log.Requesting(req.URL.String())
	response, err := client.Do(req)

	if err != nil {
		test.Log.Error(err)
		return false
	}

	// Getting the response spec.
	apiResp := test.getResponse(responseStatus, responseContentType)
	if apiResp == nil {
		test.Log.ResponseNotFound(responseContentType, responseStatus)
		return false
	}

	//	Testing the response against the spec.
	tResp := NewResponse(apiResp, test.Log)
	return tResp.Test(response)
}

func (test Operation) pickExample(examples api.ExampleList) ([]byte, string) {
	for specReqExampleName, specReqExample := range examples {
		jsonReqExample, jsonReqExampleErr := json.Marshal(specReqExample)

		if jsonReqExampleErr == nil {
			return jsonReqExample, specReqExampleName
		}

		fmt.Printf("\tThe example \"%s\" has errors: %s\n", specReqExampleName, jsonReqExampleErr.Error())
	}

	return nil, ""
}

// createRequest creates a Request instance and configures it with
// needed headers & a request body.
func (test Operation) createRequest(CT string) *http.Request {
	URL := test.createURL()

	predRequestCT := func(specCT string) bool {
		return specCT == CT
	}

	if CT == "*" {
		predRequestCT = func(specCT string) bool { return true }
	}

	// Choosing a request body.
	var reqBody *bytes.Buffer = nil

	for _, specReq := range *test.Operation.Requests {
		if predRequestCT(specReq.ContentType) {
			// Trying to find example data in the request first.
			specReqExample, specReqExampleName := test.pickExample(specReq.Examples)
			if specReqExample != nil {
				fmt.Printf("\tUsing the \"%s\" example (from operation) as request data.\n", specReqExampleName)
				reqBody = bytes.NewBuffer(specReqExample)
			} else if specReq.Schema != nil {
				// Then in the request body schema, if present.
				specReqExample, specReqExampleName := test.pickExample(specReq.Schema.Examples)
				if specReqExample != nil {
					fmt.Printf("\tUsing the \"%s\" example (from schema) as request data.\n", specReqExampleName)
					reqBody = bytes.NewBuffer(specReqExample)
				}
			}

			if reqBody != nil {
				break
			}
		}
	}

	var req *http.Request

	if reqBody != nil {
		req, _ = http.NewRequest(test.Operation.Method, URL, reqBody)
	} else {
		//TODO: check the op method name to see if body is necessary?..
		fmt.Printf("\tNo request body is available.\n")
		req, _ = http.NewRequest(test.Operation.Method, URL, nil)
	}

	req.Header.Add("Content-Type", CT)

	return req
}

// createURL creates a fully qualified URL by joining
// the server host name with an operation path
// and replaceing path parameters with actual values from `example`.
func (test Operation) createURL() string {
	path := test.Operation.Path.Path

	useParameters := func(specParams []api.Parameter, container string) {
		for _, specP := range specParams {
			if specP.In != api.ParameterLocationPath {
				continue
			}

			RX, _ := regexp.Compile("\\{" + specP.Name + "\\}")

			if RX.Match([]byte(path)) {
				if specP.Example != "" {
					path = string(RX.ReplaceAll([]byte(path), []byte(specP.Example)))
					test.Log.UsingParameterExample(&specP, container)
				} else {
					test.Log.ParameterHasNoExample(&specP, container)
				}
			}
		}
	}

	useParameters(test.Operation.Parameters, "operation")
	useParameters(test.Operation.Path.Parameters, "path")

	return fmt.Sprintf("%s%s", test.Host.URL, path)
}

// getResponse finds a response spec object to validate an actual response against.
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
