package test

import (
	"net/http"

	"github.com/x1n13y84issmd42/goasis/src/api"
	"github.com/x1n13y84issmd42/goasis/src/log"
)

// IResponse is an interface for HTTP response testers.
type IResponse interface {
	Test(resp *http.Response) bool
}

// IJSONResponse is an interface to the parser of "top level" JSON data structure.
// Arrays, Objects, as well as primitive types must be unmarshalled and validated differently,
// so there is a separate unmarshaller/parser for each of those types.
// Aggregate types use api.Schema to validate their contents against.
type IJSONResponse interface{}

// HTTPResponse --
type HTTPResponse struct {
	Log         log.ILogger
	APIResponse *api.Response
}

// NewResponse creates a new response test object. It can be either a basic HTTP response,
// or one of the JSON family of response testers.
func NewResponse(apiResp *api.Response, logger log.ILogger) IResponse {
	if apiResp.ContentType == "application/json" {
		switch apiResp.Schema.DataType {
		case api.DataTypeArray:
			return JSONResponseArray{
				Log:         logger,
				APIResponse: apiResp,
			}

		case api.DataTypeObject:
			return JSONResponseObject{
				Log:         logger,
				APIResponse: apiResp,
			}

		case api.DataTypeString:
			return JSONResponseString{
				Log:         logger,
				APIResponse: apiResp,
			}

		case api.DataTypeNumber:
			return JSONResponseNumber{
				Log:         logger,
				APIResponse: apiResp,
			}

		case api.DataTypeBoolean:
			return JSONResponseBoolean{
				Log:         logger,
				APIResponse: apiResp,
			}
		}
	}

	return HTTPResponse{
		Log:         logger,
		APIResponse: apiResp,
	}
}

// Test checks basic response properties such as status code and headers.
func (test HTTPResponse) Test(resp *http.Response) bool {
	statusOK := test.APIResponse.StatusCode == resp.StatusCode
	CTOK := (test.APIResponse.ContentType == "") || (test.APIResponse.ContentType == resp.Header.Get("Content-Type"))

	headersOK := true
	for apiHeaderName, apiHeaderValues := range test.APIResponse.Headers {
		headersOK = headersOK && NewResponseHeader(apiHeaderValues, test.Log).Test(resp.Header.Values(apiHeaderName))
	}

	if !statusOK {
		test.Log.ResponseHasWrongStatus(test.APIResponse, resp.StatusCode)
	}

	if !CTOK {
		test.Log.ResponseHasWrongContentType(test.APIResponse, resp.Header.Get("Content-Type"))
	}

	return statusOK && CTOK && headersOK
}
