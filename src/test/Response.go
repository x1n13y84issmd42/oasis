package test

import (
	"net/http"
	"strings"

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

// HTTPResponse tests basic response properties, such as status code & headers.
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

	//FIXME: this is to get rid of the "; charset=utf-8" part.
	respCT := strings.Split(resp.Header.Get("Content-Type"), ";")[0]
	CTOK := (test.APIResponse.ContentType == "") || (test.APIResponse.ContentType == respCT)

	headersOK := true
	for specHeaderName, specHeaderValues := range test.APIResponse.Headers {
		headersOK = NewResponseHeader(specHeaderValues, test.Log).Test(resp.Header.Values(specHeaderName)) && headersOK
	}

	if !statusOK {
		test.Log.ResponseHasWrongStatus(test.APIResponse, resp.StatusCode)
	}

	if !CTOK {
		test.Log.ResponseHasWrongContentType(test.APIResponse, resp.Header.Get("Content-Type"))
	}

	return statusOK && CTOK && headersOK
}
