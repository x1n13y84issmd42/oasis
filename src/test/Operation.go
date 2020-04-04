package test

import (
	"net/http"

	"github.com/x1n13y84issmd42/oasis/src/api"
	"github.com/x1n13y84issmd42/oasis/src/log"
	// "github.com/x1n13y84issmd42/oasis/src/test/security"
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

	req := specReq.CreateRequest(specHost)

	logger.Requesting(req.URL.String())
	response, err := client.Do(req)

	if err != nil {
		logger.Error(err)
		return false
	}

	// fmt.Printf("Op response: %#v\n", response)

	return Response(response, specResp, logger)
}

// SelectRequest ...
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

// SelectResponse ...
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

/* func (test Operation) pickExample(examples spec.ExampleList) ([]byte, string) {
	for specReqExampleName, specReqExample := range examples {
		jsonReqExample, jsonReqExampleErr := json.Marshal(specReqExample)

		if jsonReqExampleErr == nil {
			fmt.Printf("\tThe example \"%s\" value: '%s'\n", specReqExampleName, jsonReqExample)
			return jsonReqExample, specReqExampleName
		}

		fmt.Printf("\tThe example \"%s\" has errors: %s\n", specReqExampleName, jsonReqExampleErr.Error())
	}

	return nil, ""
} */

/* // createRequest creates a Request instance and configures it with
// needed headers, params & a request body.
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

	test.addQueryParameters(req)
	test.addHeaders(req)

	return req
} */

/* // getResponse finds a response spec object to validate an actual response against.
func (test Operation) getResponse(status int, CT string) *api.Response {
	filterCT := func(apiResp spec.Response) bool {
		return apiResp.ContentType == CT
	}

	filterStatus := func(apiResp spec.Response) bool {
		return apiResp.StatusCode == status
	}

	if CT == "*" {
		filterCT = func(spec.Response) bool { return true }
	}

	if status == 0 {
		filterStatus = func(spec.Response) bool { return true }
	}

	for _, resp := range *test.Operation.Responses {
		if filterStatus(resp) && filterCT(resp) {
			return &resp
		}
	}

	return nil
} */
